package port

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port/storage"
	"github.com/a-urth/abyr/src/service/port/storage/postgres"
)

// Service is a port service container
type Service struct {
	storage storage.Storer
}

// NewService creates and returns port service instance
func NewService() (portpb.PortServiceServer, error) {
	storage, err := postgres.NewStorer()
	if err != nil {
		return nil, err
	}

	return &Service{storage}, nil
}

// GetPort return port info for given port id
func (s *Service) GetPort(
	ctx context.Context, req *portpb.PortID,
) (*portpb.Port, error) {
	port, err := s.storage.GetPort(ctx, req.Id)
	if err != nil {
		log.WithFields(
			log.Fields{
				"service": "Port",
				"portID":  req.Id,
			},
		).Error(err)

		return nil, err
	}

	return port, nil
}

// UpsertPort upserts port entity in database from given information
func (s *Service) UpsertPort(
	ctx context.Context, req *portpb.Port,
) (*empty.Empty, error) {
	err := s.storage.UpsertPort(ctx, req)
	if err != nil {
		log.WithFields(
			log.Fields{
				"service": "Port",
				"port":    req,
			},
		).Error(err)

		return nil, err
	}

	return new(empty.Empty), nil
}
