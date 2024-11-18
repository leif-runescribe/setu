# Chat Server README

## Overview
This is a simple chat server implemented in Go, which allows clients to connect to a central server and communicate with each other. The server handles client connections, allows users to send messages, request to connect with other clients, and manage chat sessions.

### Features:
- Clients can connect to the server via TCP.
- Clients can list all currently connected clients.
- Clients can send connection requests to other clients.
- Clients can accept connection requests and start a private chat.
- Messages are transmitted between clients over the server.

## Prerequisites
To run and demo this application, you'll need:
- **Go** (version 1.18 or higher)
- **Netcat** (`nc`) or another terminal-based tool to connect to the server.

## Installation & Setup

1. **Clone the Repository:**
   If you haven’t already, clone the repository containing the code.

   ```bash
   git clone <repository-url>
   cd <repository-directory>


Build the Go Server: Build the Go application.

bash
Copy code
go build main.go

Run the Server: Start the server by running the compiled Go program.

./main

Create a client and connect to server.
nc localhost 8080
