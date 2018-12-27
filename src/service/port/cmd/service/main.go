package main

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port"
)

func main() {
	addr := ":14000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	service, err := port.NewService()
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	grpcServer := grpc.NewServer()
	portpb.RegisterPortServiceServer(grpcServer, service)

	log.Debugf("Starting to serve port service on %s", addr)

	grpcServer.Serve(lis)
}
