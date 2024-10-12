package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Client represents a connected chat user.
type Client struct {
	Conn   net.Conn
	Nick   string
	Room   *Room
	Server *Server
	writer *bufio.Writer
}

// NewClient creates a new Client from a connection.
func NewClient(conn net.Conn, server *Server) *Client {
	return &Client{
		Conn:   conn,
		Nick:   fmt.Sprintf("user_%d", time.Now().UnixNano()%10000),
		Server: server,
		writer: bufio.NewWriter(conn),
	}
}

// Send writes a message to the client's connection.
func (c *Client) Send(msg string) {
	fmt.Fprintf(c.writer, "%s\r\n", msg)
	c.writer.Flush()
}

// Handle reads input from the client and processes commands/messages.
func (c *Client) Handle() {
	defer func() {
		if c.Room != nil {
			c.Room.Leave(c)
		}
		c.Server.RemoveClient(c)
		c.Conn.Close()
	}()

	c.Send(fmt.Sprintf("Welcome to the chat server! Your nick is %s", c.Nick))
	c.Send("Commands: /nick <name>, /rooms, /join <room>, /list, /quit")
	c.Send("")

	// Auto-join the general room
	generalRoom := c.Server.GetOrCreateRoom("general")
	c.Room = generalRoom
	generalRoom.Join(c)

	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/") {
			c.handleCommand(line)
		} else {
			c.handleMessage(line)
		}
	}
}

func (c *Client) handleCommand(line string) {
	parts := strings.SplitN(line, " ", 2)
	cmd := strings.ToLower(parts[0])
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(parts[1])
	}

	switch cmd {
	case "/nick":
		if arg == "" {
			c.Send("Usage: /nick <name>")
			return
		}
		if len(arg) > 20 {
			c.Send("Error: nick must be 20 characters or less")
			return
		}
		if strings.ContainsAny(arg, " \t\r\n") {
			c.Send("Error: nick cannot contain spaces")
			return
		}
		oldNick := c.Nick
		c.Nick = arg
		c.Send(fmt.Sprintf("Nick changed to %s", c.Nick))
		if c.Room != nil {
			c.Room.Broadcast(fmt.Sprintf("*** %s is now known as %s ***", oldNick, c.Nick), c)
		}

	case "/rooms":
		rooms := c.Server.ListRooms()
		if len(rooms) == 0 {
			c.Send("No active rooms.")
			return
		}
		c.Send("Active rooms:")
		for _, info := range rooms {
			c.Send(fmt.Sprintf("  #%-20s (%d users)", info.Name, info.Count))
		}

	case "/join":
		if arg == "" {
			c.Send("Usage: /join <room>")
			return
		}
		if c.Room != nil {
			if c.Room.Name == arg {
				c.Send(fmt.Sprintf("You are already in #%s", arg))
				return
			}
			c.Room.Leave(c)
		}
		newRoom := c.Server.GetOrCreateRoom(arg)
		c.Room = newRoom
		newRoom.Join(c)

	case "/list":
		if c.Room == nil {
			c.Send("You are not in a room. Use /join <room>")
			return
		}
		members := c.Room.ListMembers()
		c.Send(fmt.Sprintf("Users in #%s (%d):", c.Room.Name, len(members)))
		for _, nick := range members {
			marker := ""
			if nick == c.Nick {
				marker = " (you)"
			}
			c.Send(fmt.Sprintf("  - %s%s", nick, marker))
		}

	case "/quit":
		c.Send("Goodbye!")
		c.Conn.Close()

	default:
		c.Send(fmt.Sprintf("Unknown command: %s", cmd))
		c.Send("Commands: /nick <name>, /rooms, /join <room>, /list, /quit")
	}
}

func (c *Client) handleMessage(msg string) {
	if c.Room == nil {
		c.Send("You are not in a room. Use /join <room>")
		return
	}

	formatted := fmt.Sprintf("[#%s] %s: %s", c.Room.Name, c.Nick, msg)
	c.Room.Broadcast(formatted, c)
}
