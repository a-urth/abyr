package storage

import (
	"context"

	"github.com/a-urth/abyr/pb/portpb"
)

// Storer defines port service's storer interface
type Storer interface {
	GetPort(ctx context.Context, portID string) (*portpb.Port, error)
	UpsertPort(ctx context.Context, port *portpb.Port) error
	DeletePort(ctx context.Context, portID string) error
	Close() error
}
