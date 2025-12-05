package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)


type Message struct {
	UserID string
	Text   string
}


type ChatServer struct {
	clients  map[string]chan string
	history  []string
	mu       sync.Mutex
}


func (s *ChatServer) Join(userID string, ack *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.clients == nil {
		s.clients = make(map[string]chan string)
	}

	ch := make(chan string, 50) 
	s.clients[userID] = ch

	joinMsg := fmt.Sprintf("User [%s] joined", userID)
	for id, clientCh := range s.clients {
		if id != userID {
			clientCh <- joinMsg
		}
	}

	
	for _, msg := range s.history {
		ch <- msg
	}

	*ack = "Joined successfully"
	fmt.Println(joinMsg)
	return nil
}


func (s *ChatServer) SendMessage(msg Message, ack *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fullMsg := fmt.Sprintf("[%s]: %s", msg.UserID, msg.Text)
	s.history = append(s.history, fullMsg)

	for id, clientCh := range s.clients {
		if id != msg.UserID { 
			clientCh <- fullMsg
		}
	}

	*ack = "Message sent"
	fmt.Println(fullMsg)
	return nil
}


func (s *ChatServer) Receive(userID string, msg *string) error {
	s.mu.Lock()
	ch, ok := s.clients[userID]
	s.mu.Unlock()

	if !ok {
		return fmt.Errorf("client not found")
	}

	*msg = <-ch
	return nil
}

func main() {
	server := new(ChatServer)
	rpc.Register(server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
