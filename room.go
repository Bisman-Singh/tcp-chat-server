package main

import (
	"fmt"
	"sync"
)

// Room represents a chat room that holds clients and broadcasts messages.
type Room struct {
	Name    string
	clients map[*Client]bool
	mu      sync.RWMutex
}

// NewRoom creates a new Room with the given name.
func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		clients: make(map[*Client]bool),
	}
}

// Join adds a client to the room and announces the join.
func (r *Room) Join(c *Client) {
	r.mu.Lock()
	r.clients[c] = true
	r.mu.Unlock()

	r.Broadcast(fmt.Sprintf("*** %s has joined #%s ***", c.Nick, r.Name), c)
	c.Send(fmt.Sprintf("Joined room #%s (%d user(s) online)", r.Name, r.Count()))
}

// Leave removes a client from the room and announces the departure.
func (r *Room) Leave(c *Client) {
	r.mu.Lock()
	delete(r.clients, c)
	r.mu.Unlock()

	r.Broadcast(fmt.Sprintf("*** %s has left #%s ***", c.Nick, r.Name), nil)
}

// Broadcast sends a message to all clients in the room, optionally excluding one.
func (r *Room) Broadcast(msg string, exclude *Client) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		if client != exclude {
			client.Send(msg)
		}
	}
}

// Count returns the number of clients in the room.
func (r *Room) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

// ListMembers returns a list of nicknames in the room.
func (r *Room) ListMembers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	members := make([]string, 0, len(r.clients))
	for client := range r.clients {
		members = append(members, client.Nick)
	}
	return members
}
