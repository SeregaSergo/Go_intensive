package main

import (
	"fmt"
	pb "github.com/SeregaSergo/Go_intensive/team_0/api"
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

	s := pb.Transmitter{}

	grpcServer := grpc.NewServer()

	pb.RegisterCommunicationServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
