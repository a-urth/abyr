package clientapi

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/a-urth/abyr/pb/portpb"
	log "github.com/sirupsen/logrus"
)

func portsReader(
	ctx context.Context,
	ports chan<- portpb.Port,
	wg *sync.WaitGroup,
	filePath string,
) {
	log.Debug("ports reader started")

	wg.Add(1)
	defer func() {
		wg.Done()
		log.Debug("ports reader finished")
	}()

	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("error during file opening - %v", err)
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	// skip opening brace
	_, err = dec.Token()
	if err != nil {
		log.Errorf("error during file opening - %v", err)
		return
	}

	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				log.Errorf("error during parsing - %v", err)
			}
			break
		}

		var port portpb.Port
		if err := dec.Decode(&port); err != nil {
			if err != io.EOF {
				log.Errorf("error during parsing - %v", err)
			}

			break
		}

		port.Id = t.(string)
		ports <- port
	}
}

func portsSender(
	ctx context.Context,
	ports <-chan portpb.Port,
	wg *sync.WaitGroup,
	portClient portpb.PortServiceClient,
) {
	log.Debug("ports sender started")

	wg.Add(1)
	defer func() {
		wg.Done()
		log.Debug("ports sender finished")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case port := <-ports:
			if _, err := portClient.UpsertPort(ctx, &port); err != nil {
				log.WithFields(log.Fields{
					"port": port,
				}).Error(err)
			}
		}
	}
}
