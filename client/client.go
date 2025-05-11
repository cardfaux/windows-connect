// client.go
package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.0.151:4444")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")
	conn.Write([]byte("Hello from client"))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	fmt.Printf("Received from server: %s\n", string(buffer[:n]))
}
