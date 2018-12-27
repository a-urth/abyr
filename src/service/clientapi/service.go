package clientapi

import (
	"context"

	"github.com/a-urth/abyr/pb/clientapipb"
)

// Service represents client api service structure
type Service struct {
	FilePath  string
	QueueSize int
}

// GetPort retrieves port information for given port id
func (s *Service) GetPort(
	ctx context.Context, req *clientapipb.PortID,
) (*clientapipb.Port, error) {
	return nil, nil
}
