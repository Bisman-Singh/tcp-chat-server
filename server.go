package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// RoomInfo holds summary info about a room.
type RoomInfo struct {
	Name  string
	Count int
}

// Server manages all rooms and connected clients.
type Server struct {
	rooms   map[string]*Room
	clients map[*Client]bool
	mu      sync.RWMutex
}

// NewServer creates a new Server.
func NewServer() *Server {
	s := &Server{
		rooms:   make(map[string]*Room),
		clients: make(map[*Client]bool),
	}
	// Create default room
	s.rooms["general"] = NewRoom("general")
	return s
}

// Start listens on the given address and accepts connections.
func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	defer listener.Close()

	log.Printf("Chat server listening on %s", addr)
	log.Printf("Default room: #general")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		client := NewClient(conn, s)
		s.addClient(client)

		log.Printf("New connection from %s (nick: %s)", conn.RemoteAddr(), client.Nick)

		go client.Handle()
	}
}

func (s *Server) addClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c] = true
}

// RemoveClient removes a client from the server tracking.
func (s *Server) RemoveClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c)
	log.Printf("Client disconnected: %s (%s)", c.Nick, c.Conn.RemoteAddr())
}

// GetOrCreateRoom returns an existing room or creates a new one.
func (s *Server) GetOrCreateRoom(name string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.rooms[name]; exists {
		return room
	}

	room := NewRoom(name)
	s.rooms[name] = room
	log.Printf("Room created: #%s", name)
	return room
}

// ListRooms returns a list of all rooms with their member counts.
func (s *Server) ListRooms() []RoomInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var infos []RoomInfo
	for _, room := range s.rooms {
		infos = append(infos, RoomInfo{
			Name:  room.Name,
			Count: room.Count(),
		})
	}
	return infos
}
