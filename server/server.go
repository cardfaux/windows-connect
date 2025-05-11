// server.go
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

func (s *server) Echo(ctx context.Context, req *grpcapi.EchoRequest) (*grpcapi.EchoResponse, error) {
	log.Printf("Received: %s\n", req.GetMessage())
	return &grpcapi.EchoResponse{Reply: "Hello from server"}, nil
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