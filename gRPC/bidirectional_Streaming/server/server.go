package main

import (
	"fmt"
	"log"
	"net"

	pb "Bidirectional/bidirectional"

	"google.golang.org/grpc"
)

// BidirectionalService 서버 구조체 정의
type BidirectionalServer struct {
	pb.UnimplementedBidirectionalServer
}

// GetServerResponse 메서드는 양방향 스트리밍 처리를 수행합니다.
func (s *BidirectionalServer) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {
	fmt.Println("Server processing gRPC bidirectional streaming.")

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Client disconnected or error: %v", err)
			return err
		}

		fmt.Printf("[client to server] %s\n", msg.Message)

		// 클라이언트로 받은 메시지를 다시 반환합니다.
		if err := stream.Send(msg); err != nil {
			log.Printf("Failed to send message back to client: %v", err)
			return err
		}
	}
}

func serve() {
	// gRPC 서버 초기화
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterBidirectionalServer(server, &BidirectionalServer{})

	fmt.Println("Starting server. Listening on port 50051.")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	serve()
}
