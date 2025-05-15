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

# Run This With Different Shells From The Server

## Running Commands on the Client: Shell Prefix Usage

When sending commands from the server to the client, you can specify which shell the client should use to execute the command by **prefixing the command string** with the shell name followed by a colon (`:`).

### Supported Shell Prefixes and Examples

| Platform        | Shell        | Example Command                |
| --------------- | ------------ | ------------------------------ |
| **Windows**     | `cmd`        | `cmd: dir /w`                  |
| **Windows**     | `powershell` | `powershell: Get-Process`      |
| **macOS/Linux** | `sh`         | `sh: ls -la`                   |
| **macOS/Linux** | `bash`       | `bash: echo "Hello from bash"` |

### How It Works

- The server sends the full command string including the prefix to the client.
- The client reads the prefix (e.g., `cmd`, `powershell`, `sh`, or `bash`) and runs the rest of the command using the specified shell.
- If no prefix is provided, the client can default to a predefined shell depending on the operating system.

---

### Example Usage

To list files on a Windows client using `cmd`:

```text
cmd: dir /w
```

To list processes on a Windows client using PowerShell:

```text
powershell: Get-Process
```

To list files on a Linux or macOS client using `sh`:

```text
sh: ls -la
```

===Prefixes are case-insensitive (e.g., CMD:, cmd:, or Cmd: all work).===

# What Is gRPC Streaming?

```lua
╭───────────────────────────── SERVER ─────────────────────────────╮
│                                                                 │
│  [stdin]         [Send command to client]                       │
│     │                     │                                     │
│     ▼                     ▼                                     │
│  User types ───▶ stream.Send(CommandResponse{output: cmd})      │
│                                                                 │
│                           │                                     │
│                           ▼                                     │
│             [Receive output from client]                        │
│                    stream.Recv()                               │
│                           │                                     │
│                           ▼                                     │
│                  Print output to terminal                       │
╰─────────────────────────────────────────────────────────────────╯

                ▲                                         │
                │                                         ▼

╭────────────────────────────── CLIENT ─────────────────────────────╮
│                                                                   │
│       [Receive command from server]                               │
│               stream.Recv()                                       │
│                     │                                             │
│                     ▼                                             │
│   Parse "shell:command" prefix (e.g., bash: ls -l)                │
│                     │                                             │
│                     ▼                                             │
│   Run exec.Command(shell, args)                                   │
│                     │                                             │
│                     ▼                                             │
│     Capture stdout/stderr output                                  │
│                     │                                             │
│                     ▼                                             │
│ stream.Send(CommandRequest{command: output, shell: usedShell})   │
╰───────────────────────────────────────────────────────────────────╯
```

### 🔄 Key Flow

- 🔼 Server → Client: CommandResponse.output = command string
- 🔽 Client → Server: CommandRequest.command = output, CommandRequest.shell = shell used

# HOW TO RUN THE COMMANDS FOR THE LATEST VERSION

## How to Run Commands on the gRPC Server

Once your gRPC server is running and connected to a client, you can enter commands directly into the server console. These commands tell the client what to do, such as listing files, reading a file, or getting file info.

---

### Command Format

#### `COMMAND_TYPE:ARGUMENT`

- `COMMAND_TYPE` — The type of command you want to execute. Currently supported commands are:

  - `LIST_FILES` — List all files and directories inside the specified directory.
  - `GET_FILE` — Retrieve the contents of a specified file.
  - `GET_INFO` — Get information about a file or directory (size, modification date, etc.).

- `ARGUMENT` — The path or target for the command (e.g., a directory path for `LIST_FILES` or a file path for `GET_FILE`).

---

## Examples

### List files in a directory

To list all files and folders inside `/home/user` (Linux/macOS):

`LIST_FILES:/home/user`

On Windows, to list files in `C:\Users\YourName\Documents`:

`LIST_FILES:C:\Users\YourName\Documents`

### Read the contents of a file

To get the contents of `/etc/hosts` (Linux/macOS):

`GET_FILE:/etc/hosts`

On Windows, to read `C:\Windows\System32\drivers\etc\hosts`:

`GET_FILE:C:\Windows\System32\drivers\etc\hosts`

---

### Get information about a file or directory

To get info about `/var/log`:

`GET_INFO:/var/log`

To get info about `C:\Program Files` on Windows:

`GET_INFO:C:\Program Files`

---

## Notes

- Make sure the client machine has permission to access the specified path.
- Paths should be valid on the **client machine**, not the server.
- If a command or path is invalid, you will get an error response printed in the server console.

---

## Summary

| Command    | Example                         | Purpose                     |
| ---------- | ------------------------------- | --------------------------- |
| LIST_FILES | `LIST_FILES:/path/to/directory` | List files and directories  |
| GET_FILE   | `GET_FILE:/path/to/file`        | Retrieve file contents      |
| GET_INFO   | `GET_INFO:/path/to/file_or_dir` | Get metadata about file/dir |

---

Enter these commands in the server terminal prompt and press **Enter** to send them to the client and see the results.
