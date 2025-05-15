package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

type server struct {
	grpcapi.UnimplementedEchoServiceServer
}

func (s *server) ExecuteCommand(stream grpcapi.EchoService_ExecuteCommandServer) error {
	// Channel to send commands from stdin
	cmdChan := make(chan *grpcapi.CommandMessage)

	// Goroutine: read commands from stdin and send them as CommandRequest messages
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter command (prefix with shell if needed, e.g. 'bash: ls'): ")
			if !scanner.Scan() {
				close(cmdChan)
				return
			}
			text := scanner.Text()

			// Parse shell prefix if present
			shell := ""
			command := text
			if idx := indexOfColon(text); idx != -1 {
				shell = text[:idx]
				command = text[idx+1:]
			}

			cmdMsg := &grpcapi.CommandMessage{
				Message: &grpcapi.CommandMessage_CommandRequest{
					CommandRequest: &grpcapi.CommandRequest{
						Command: command,
						Shell:   shell,
					},
				},
			}
			cmdChan <- cmdMsg
		}
	}()

	// Goroutine: receive outputs from client
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving from client: %v", err)
				return
			}
			switch m := msg.Message.(type) {
			case *grpcapi.CommandMessage_CommandResponse:
				output := m.CommandResponse.Output
				log.Printf("Output from client:\n%s\n", output)
			default:
				log.Printf("Received unexpected message type from client")
			}
		}
	}()

	// Send commands from cmdChan to client
	for cmd := range cmdChan {
		if err := stream.Send(cmd); err != nil {
			log.Printf("Error sending command to client: %v", err)
			return err
		}
	}

	return nil
}

// indexOfColon returns the index of the first colon in s or -1 if not found
func indexOfColon(s string) int {
	for i, c := range s {
		if c == ':' {
			return i
		}
	}
	return -1
}

func main() {
	lis, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	grpcapi.RegisterEchoServiceServer(s, &server{})

	fmt.Println("gRPC server listening on :4444...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}