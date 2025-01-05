package main

import (
	"io"
	"log"
	"net"

	pb "streaming/streaming"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedClientStreamingServer
}

func (s *server) GetServerResponse(stream pb.ClientStreaming_GetServerResponseServer) error {
	log.Println("Server processing gRPC client-streaming.")
	count := int32(0)

	// Process messages from the stream
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			// End of client stream, send response
			return stream.SendAndClose(&pb.Number{Value: count})
		}
		if err != nil {
			log.Fatalf("Error reading from stream: %v", err)
			return err
		}
		count++
	}
}

func main() {
	// Start a TCP listener
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to open port: %v", err)
	}
	log.Println("Server started. Listening on port 50051.")

	// Create and start the gRPC server
	grpcServer := grpc.NewServer()

	pb.RegisterClientStreamingServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
