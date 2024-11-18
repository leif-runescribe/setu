// main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"setu/utils"
	"sync"
)

type Client struct {
	id       string
	conn     net.Conn
	outbound chan string
	paired   *Client
	server   *Server
}

type Server struct {
	clients  map[string]*Client
	mutex    sync.RWMutex
	requests map[string]string
}

func NewServer() *Server {
	return &Server{
		clients:  make(map[string]*Client),
		requests: make(map[string]string),
	}
}

func (s *Server) registerClient(conn net.Conn) *Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := utils.GenId()
	client := &Client{
		id:       id,
		conn:     conn,
		outbound: make(chan string, 10),
		server:   s,
	}
	s.clients[id] = client

	return client
}

func (s *Server) removeClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client.paired != nil {
		client.paired.paired = nil
		client.paired.outbound <- "Your chat partner has disconnected.\n"
	}
	delete(s.clients, client.id)
	delete(s.requests, client.id)
}

func (s *Server) listClients() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var clientList []string
	for id := range s.clients {
		clientList = append(clientList, id)
	}
	return clientList
}

func (c *Client) handleConnection() {
	defer func() {

		c.server.removeClient(c)
		c.conn.Close()
	}()

	c.outbound <- fmt.Sprintf("Welcome! Your ID is %s\nCommands:\n/list - Show all clients\n/connect <client-id> - Request to connect\n/accept - Accept pending request\n/quit - Disconnect\n", c.id)

	go c.handleOutbound()

	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		msg := scanner.Text()

		if msg == "/quit" {
			return
		}

		if !c.handleCommand(msg) && c.paired != nil {
			c.paired.outbound <- fmt.Sprintf("%s: %s", c.id, msg)
		}
	}
}

func (c *Client) handleOutbound() {
	for msg := range c.outbound {
		fmt.Fprintf(c.conn, "%s\n", msg)
	}
}

func (c *Client) handleCommand(msg string) bool {
	if msg == "/list" {
		clients := c.server.listClients()
		c.outbound <- fmt.Sprintf("Connected clients: %v\n", clients)
		return true
	}

	var targetID string
	n, _ := fmt.Sscanf(msg, "/connect %s", &targetID)
	if n == 1 {
		c.requestConnection(targetID)
		return true
	}

	if msg == "/accept" {
		c.acceptConnection()
		return true
	}

	return false
}

func (c *Client) requestConnection(targetID string) {
	c.server.mutex.Lock()
	defer c.server.mutex.Unlock()

	target, exists := c.server.clients[targetID]
	if !exists {
		c.outbound <- "Client not found"
		return
	}

	if target.paired != nil {
		c.outbound <- "Client is already in a chat"
		return
	}

	if c.paired != nil {
		c.outbound <- "You are already in a chat"
		return
	}

	c.server.requests[targetID] = c.id
	c.outbound <- fmt.Sprintf("Connection request sent to %s\n", targetID)
	target.outbound <- fmt.Sprintf("%s wants to chat with you. Type /accept to accept.\n", c.id)
}

func (c *Client) acceptConnection() {
	c.server.mutex.Lock()
	defer c.server.mutex.Unlock()

	requestorID, exists := c.server.requests[c.id]
	if !exists {
		c.outbound <- "No pending chat requests"
		return
	}

	requestor, exists := c.server.clients[requestorID]
	if !exists {
		c.outbound <- "Requesting client no longer connected"
		delete(c.server.requests, c.id)
		return
	}

	c.paired = requestor
	requestor.paired = c
	delete(c.server.requests, c.id)

	c.outbound <- fmt.Sprintf("You are now chatting with %s", requestorID)
	requestor.outbound <- fmt.Sprintf("You are now chatting with %s", c.id)
}

func main() {
	server := NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		listener.Close()
	}()
	fmt.Println("Chat server running on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		client := server.registerClient(conn)
		fmt.Printf("New client registered - %s\n", client.id)
		go client.handleConnection()
	}
}
