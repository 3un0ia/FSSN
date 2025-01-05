package main

import (
	"context"
	"log"
	"net"

	pb "grpc_grpc/hello_grpc"

	grpc "google.golang.org/grpc"
)

// (4) MyServiceServer는 gRPC 서버 구조체
type MyServiceServer struct {
	pb.UnimplementedMyServiceServer // gRPC 기본 구현 포함
}

// (5.1) proto 파일 내 정의된 RPC 함수에 대응하는 메서드 작성
func (s *MyServiceServer) MyFunction(ctx context.Context, req *pb.MyNumber) (*pb.MyNumber, error) {
	// 응답 메시지 생성 및 값 설정
	result := pb.My_func(req.Value)

	response := &pb.MyNumber{Value: result}
	return response, nil
}

func main() {
	// (6) gRPC 서버의 포트를 열기
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// (7) gRPC 서버 생성
	server := grpc.NewServer()

	// (7.1) MyServiceServer 등록
	pb.RegisterMyServiceServer(server, &MyServiceServer{})

	// (8) 서버 실행
	log.Println("Starting server. Listening on port 50051.")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
