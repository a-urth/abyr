package client

import (
	"io"

	"google.golang.org/grpc"

	"github.com/a-urth/abyr/pb/portpb"
)

// Client creates and returns client for port service
func Client(serverAddr string) (portpb.PortServiceClient, io.Closer, error) {
	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithInsecure(), // assume that we are running all services in a secured network
	)
	if err != nil {
		return nil, nil, err
	}

	return portpb.NewPortServiceClient(conn), conn, nil
}
