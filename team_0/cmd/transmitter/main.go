package main

import (
	"fmt"
	"github.com/SeregaSergo/Go_intensive/team_0/internal/transport"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Println("Transmitter is started!")

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := transport.Transmitter{}

	grpcServer := grpc.NewServer()

	transport.RegisterCommunicationServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
