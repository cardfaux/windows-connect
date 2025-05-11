// server.go
package main

import (
	"fmt"
	"net"
	//"github.com/cardfaux/windows-connect/grpcapi"
)

func main() {
	ln, err := net.Listen("tcp", ":4444") // listen on all interfaces
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is listening on port 4444...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	fmt.Printf("Received from client: %s\n", string(buffer[:n]))
	conn.Write([]byte("Hello from server"))
}
