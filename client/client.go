package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cardfaux/windows-connect/grpcapi"
	"google.golang.org/grpc"
)

func main() {
    //conn, err := grpc.Dial("192.168.0.151:4444", grpc.WithInsecure())
		conn, err := grpc.Dial("6.tcp.ngrok.io:17905", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := grpcapi.NewEchoServiceClient(conn)

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Enter message: ")
        if !scanner.Scan() {
            break
        }
        msg := scanner.Text()

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        resp, err := client.Echo(ctx, &grpcapi.EchoRequest{Message: msg})
        if err != nil {
            log.Printf("Error calling Echo: %v", err)
            continue
        }
        fmt.Printf("Response from server: %s\n", resp.GetReply())
    }
}