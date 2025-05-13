package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

type server struct {
	grpcapi.UnimplementedEchoServiceServer
}

func (s *server) ExecuteCommand(ctx context.Context, req *grpcapi.CommandRequest) (*grpcapi.CommandResponse, error) {
	log.Printf("Client asked to execute: %s", req.Command)
	return &grpcapi.CommandResponse{
		Output: "This should be implemented on the client, not server.",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	grpcapi.RegisterEchoServiceServer(s, &server{})

	fmt.Println("gRPC server is listening on port 4444...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}