// client.go (TCP version with persistent connection)
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	//"github.com/cardfaux/windows-connect/grpcapi"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.0.151:4444")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")
	conn.Write([]byte("Hello from client"))

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			fmt.Printf("Received from server: %s\n", string(buffer[:n]))
		}
	}()

	// Allow user to send messages to server
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		if scanner.Scan() {
			msg := scanner.Text()
			conn.Write([]byte(msg))
		} else {
			break
		}
	}
}
