package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "streaming/streaming"

	"google.golang.org/grpc"
)

// generateMessages generates a series of messages to send to the server.
func generateMessages() []*pb.Message {
	return []*pb.Message{
		{Message: "message #1"},
		{Message: "message #2"},
		{Message: "message #3"},
		{Message: "message #4"},
		{Message: "message #5"},
	}
}

// sendMessage establishes a streaming RPC and sends messages to the server.
func sendMessage(client pb.ClientStreamingClient) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Open a stream
	stream, err := client.GetServerResponse(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Send messages
	messages := generateMessages()
	for _, msg := range messages {
		log.Printf("[client to server] %s", msg.Message)
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
		time.Sleep(500 * time.Millisecond) // Simulate delay between messages
	}

	// Close the stream for sending messages
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close sending stream: %v", err)
	}

	// Receive response from the server
	response, err := stream.CloseAndRecv()
	if err == io.EOF {
		log.Println("Server closed the connection.")
	} else if err != nil {
		log.Fatalf("Error receiving server response: %v", err)
	} else {
		log.Printf("[server to client] %d", response.Value)
	}
}

func main() {
	// Create a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a client stub
	client := pb.NewClientStreamingClient(conn)

	// Send messages using the stub
	sendMessage(client)
}
