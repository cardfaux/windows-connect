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
	log.Printf("Server received command: %s", req.Command)

	// Send the response back to the client (you can modify this to fit your needs)
	return &grpcapi.CommandResponse{
		Output: "Command received by client: " + req.Command,
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
