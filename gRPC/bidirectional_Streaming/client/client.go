package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "Bidirectional/bidirectional"
)

func makeMessage(message string) *pb.Message {
	return &pb.Message{
		Message: message,
	}
}

func generateMessages() []*pb.Message {
	messages := []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
	return messages
}

func sendMessage(stub pb.BidirectionalClient) {
	// Context 생성 (타임아웃 포함 가능)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := stub.GetServerResponse(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	go func() {
		for _, msg := range generateMessages() {
			fmt.Printf("[client to server] %s\n", msg.Message)
			if err := stream.Send(msg); err != nil {
				log.Fatalf("Failed to send message: %v", err)
				return
			}
		}
	}()

	for {
		serverMsg, err := stream.Recv()
		if err != nil {
			log.Printf("Connection closed: %v", err)
			break
		}

		fmt.Printf("[server to client] %s\n", serverMsg.Message)
	}
}

func run() {
	channel, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer channel.Close()

	stub := pb.NewBidirectionalClient(channel)

	sendMessage(stub)
}

func main() {
	run()
}
