package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct {
	grpcapi.UnimplementedEchoServiceServer
}

func (s *server) ExecuteCommand(stream grpcapi.EchoService_ExecuteCommandServer) error {
	done := make(chan struct{})
	var once sync.Once

	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	// Receive responses concurrently
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				closeDone()
				return
			}
			if resp := msg.GetCommandResponse(); resp != nil {
				fmt.Printf("\nResponse:\n%s\n\n", resp.Output)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-done:
			log.Println("Client disconnected or receive failed — stopping input loop.")
			return nil
		default:
			fmt.Print("Enter command (e.g., LIST_FILES:/ or GET_FILE:/path/to/file): ")
			if !scanner.Scan() {
				log.Println("Input scan failed or EOF")
				closeDone()
				return nil
			}

			text := scanner.Text()
			parts := strings.SplitN(text, ":", 2)
			if len(parts) == 0 {
				fmt.Println("Invalid command format.")
				continue
			}

			cmdStr := parts[0]
			arg := ""
			if len(parts) > 1 {
				arg = parts[1]
			}

			var cmdEnum grpcapi.CommandType
			switch strings.ToUpper(cmdStr) {
			case "LIST_FILES":
				cmdEnum = grpcapi.CommandType_LIST_FILES
			case "GET_FILE":
				cmdEnum = grpcapi.CommandType_GET_FILE
			case "GET_INFO":
				cmdEnum = grpcapi.CommandType_GET_INFO
			default:
				fmt.Println("Unknown command. Valid commands: LIST_FILES, GET_FILE, GET_INFO")
				continue
			}

			err := stream.Send(&grpcapi.CommandMessage{
				Message: &grpcapi.CommandMessage_CommandRequest{
					CommandRequest: &grpcapi.CommandRequest{
						Command:  cmdEnum,
						Argument: arg,
					},
				},
			})
			if err != nil {
				log.Printf("Send error: %v", err)
				closeDone()
				return err
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	)

	grpcapi.RegisterEchoServiceServer(grpcServer, &server{})

	fmt.Println("gRPC server listening on :4444")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
