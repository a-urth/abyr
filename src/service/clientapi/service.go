package clientapi

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port/client"
)

// Service represents client api service structure
type Service struct {
	portClient       portpb.PortServiceClient
	portClientCloser io.Closer

	queue chan portpb.Port

	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// Servicer defines client api service interface
// its a stupid name, but cannot come up with something better atm
type Servicer interface {
	Ping(w http.ResponseWriter, r *http.Request)
	GetPort(w http.ResponseWriter, r *http.Request)
	Close()
}

// NewService creates and returns port service instance
// TODO: proper configuration should be used here
func NewService() (Servicer, error) {
	portClient, closer, err := client.NewClient(":14000")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	portsChan := make(chan portpb.Port, 50)

	service := Service{
		portClient:       portClient,
		portClientCloser: closer,
		cancel:           cancel,
		wg:               wg,
	}

	go portsSender(ctx, portsChan, wg, portClient)
	go portsReader(ctx, portsChan, wg, "ports.json")

	return &service, nil
}

// Ping responds with dummy message to any request
func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// Close closes all resources
func (s *Service) Close() {
	// first cancel workers
	s.cancel()

	// then wait them to finish
	s.wg.Wait()

	// then close client
	if err := s.portClientCloser.Close(); err != nil {
		log.Errorf("Cannot close port service client - %v", err)
	}
}

// GetPort return port by given id
func (s *Service) GetPort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portID := vars["id"]

	// here should be proper context usage, at least with gorilla
	ctx := context.TODO()

	port, err := s.portClient.GetPort(ctx, &portpb.PortID{Id: portID})
	if err != nil {
		log.WithFields(log.Fields{
			"service": "ClientAPI", // this should be in context of logger, right
			"portID":  portID,
		}).Error(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	marshaler := jsonpb.Marshaler{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := marshaler.Marshal(w, port); err != nil {
		log.WithFields(log.Fields{
			"service": "ClientAPI",
			"portID":  portID,
		}).Error(err)
	}
}
