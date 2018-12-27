package main

import (
	"net"

	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	
	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port"
)

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	addr := ":14000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	service, dbCloser, err := port.NewService()
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	defer dbCloser.Close()

	grpcServer := grpc.NewServer()
	portpb.RegisterPortServiceServer(grpcServer, service)

	log.Debugf("Starting to serve port service on %s", addr)

	grpcServer.Serve(lis)
}
