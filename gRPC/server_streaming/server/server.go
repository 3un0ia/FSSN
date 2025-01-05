package main

import (
	"log"
	"net"

	pb "streaming/streaming"

	"google.golang.org/grpc"
)

// makeMessage creates a Message instance.
func makeMessage(message string) *pb.Message {
	return &pb.Message{
		Message: message,
	}
}

// ServerStreamingService implements the ServerStreaming gRPC service.
type ServerStreamingService struct {
	pb.UnimplementedServerStreamingServer
}

// GetServerResponse streams multiple messages back to the client.
func (s *ServerStreamingService) GetServerResponse(req *pb.Number, stream pb.ServerStreaming_GetServerResponseServer) error {
	messages := []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}

	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Create a gRPC server
	server := grpc.NewServer()
	pb.RegisterServerStreamingServer(server, &ServerStreamingService{})

	// Start listening on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting server. Listening on port 50051.")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
