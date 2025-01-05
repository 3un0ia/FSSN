package main

import (
	"context"
	"log"

	pb "streaming/streaming"

	"google.golang.org/grpc"
)

func recvMessage(client pb.ServerStreamingClient) {
	request := &pb.Number{Value: 5}

	// Calling the GetServerResponse method
	stream, err := client.GetServerResponse(context.Background(), request)
	if err != nil {
		log.Fatalf("Error calling GetServerResponse: %v", err)
	}

	// Receiving messages from the stream
	for {
		response, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("[server to client] %s", response.Message)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerStreamingClient(conn)
	recvMessage(client)
}
