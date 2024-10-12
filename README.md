# TCP Chat Server

A multi-client TCP chat server written in Go. Supports multiple rooms, nick changes, and real-time message broadcasting.

## Features

- Multi-client TCP connections
- Chat rooms with join/leave announcements
- Default "general" room for new connections
- Commands: `/nick`, `/rooms`, `/join`, `/quit`, `/list`
- Broadcast messages to all room members
- Thread-safe room and client management
- Automatic nick assignment on connect

## Build

```bash
go build -o tcp-chat-server .
```

## Usage

### Starting the server

```bash
# Start on default port 8080
./tcp-chat-server

# Start on a custom port
./tcp-chat-server -port 9000

# Bind to specific interface
./tcp-chat-server -host 127.0.0.1 -port 8080
```

### Connecting as a client

Use any TCP client such as `telnet` or `nc`:

```bash
# Using telnet
telnet localhost 8080

# Using netcat
nc localhost 8080
```

### Chat Commands

| Command          | Description                          |
|------------------|--------------------------------------|
| `/nick <name>`   | Change your nickname                 |
| `/rooms`         | List all active rooms                |
| `/join <room>`   | Join or create a room                |
| `/list`          | List users in your current room      |
| `/quit`          | Disconnect from the server           |

### Flags

| Flag    | Default   | Description           |
|---------|-----------|-----------------------|
| `-host` | `0.0.0.0` | Host to bind to       |
| `-port` | `8080`    | Port to listen on     |

### Example Session

```
Welcome to the chat server! Your nick is user_4821
Commands: /nick <name>, /rooms, /join <room>, /list, /quit

Joined room #general (1 user(s) online)
/nick alice
Nick changed to alice
Hello everyone!
*** bob has joined #general ***
[#general] bob: Hey alice!
/rooms
Active rooms:
  #general             (2 users)
/join devops
Joined room #devops (1 user(s) online)
/list
Users in #devops (1):
  - alice (you)
/quit
Goodbye!
```


