# 📡 Go Client-Server Test (Local Network Communication)

This project is a **basic client-server application** written in Go for testing network communication between two machines—typically a **macOS host** (server) and a **Windows VM** (client).

---

## 🔧 Project Structure

.
├── README.md # You're reading this!
└── Windows-Connect
├── client
│ └── client.go # Connects to the server and sends a test message
└── server
└── server.go # Listens for connections and prints received messages

## 🖥️ Server (macOS / Host Machine)

The server:

- Listens on a specified IP (e.g. `192.168.0.151`) and port (e.g. `4444`)
- Accepts TCP connections
- Prints any messages received from a client

### Usage:

```bash
go run server.go
```

## 🖥️ Client (Windows / Virtual Machine)

The client:

- Connects to the server's IP and port
- Sends a simple message like "Hello from client!"
- Closes the connection

### Update the Server IP

In `client.go`, make sure this line reflects your host machine's IP:

```bash
conn, err := net.Dial("tcp", "192.168.0.151:4444")
```

Replace 192.168.0.151 with your actual server IP.

### Build the client for Windows

To compile the client as a `.exe` for your Windows VM:

```bash
GOOS=windows GOARCH=amd64 go build -o client.exe client.go
```

Then transfer and run `client.exe` on your Windows machine.

## 🛠️ Configuration Notes

- Both machines must be on the same local network.
- The server IP in `client.go` must match the IP shown when running `ifconfig` (e.g., `192.168.0.151`).
- Ensure the server’s port (e.g., `4444`) is open and not blocked by firewalls.
- Use `netstat -an | grep 4444` to confirm the server is listening.

## ✅ Example Output

### Server

```bash
Server listening on 192.168.0.151:4444
Connection received from: 192.168.0.157
Received: Hello from client!
```

### Client

```bash
Sent message to server.
```

## 🧪 Purpose

This project exists solely for testing purposes. It helps verify basic TCP/IP connectivity between two machines using Go. There is no malicious code or logging—just a simple message-sending mechanism for connection validation.

# How To Run This Test Windows Connection With NGROK

1. Run the server with `go run server.go`
2. Expose the gRPC port (e.g., 4444) by running `ngrok tcp 4444` in the terminal
3. Copy the forwarded TCP address shown by ngrok, e.g.: `Forwarding tcp://0.tcp.ngrok.io:12345 -> localhost:4444`
4. On your client (remote machine), connect to the ngrok address: `grpc.Dial("0.tcp.ngrok.io:12345", grpc.WithInsecure())`
