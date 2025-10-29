package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// Message represents a single chat message
type Message struct {
	Sender  string
	Content string
}

// ChatServer stores all messages and handles RPC calls
type ChatServer struct {
	messages []Message
	mu       sync.Mutex
}

// SendMessage is called by the client to send a new message
func (s *ChatServer) SendMessage(msg Message, reply *[]Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store message
	s.messages = append(s.messages, msg)

	// Print to server console
	fmt.Printf("%s: %s\n", msg.Sender, msg.Content)

	// Return updated chat history
	*reply = s.messages
	return nil
}

// GetHistory returns all messages to a client
func (s *ChatServer) GetHistory(dummy int, reply *[]Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	*reply = s.messages
	return nil
}

func main() {
	server := new(ChatServer)
	err := rpc.Register(server)
	if err != nil {
		log.Fatalf("Error registering RPC server: %v", err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Listener error: %v", err)
	}

	fmt.Println("Chat server running on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
