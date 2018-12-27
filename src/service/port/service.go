package port

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port/storage"
)

// Service is a port service container
type Service struct {
	storage storage.Storer
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
				"method":  "GetPort",
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
				"method":  "UpsertPort",
				"port":    req,
			},
		).Error(err)

		return nil, err
	}

	return new(empty.Empty), nil
}
