package main

import (
	"context"
	"fmt"
	"time"

	pb "grpc_grpc/hello_grpc"

	grpc "google.golang.org/grpc"
)

func main() {
	// (3) gRPC 통신 채널을 생성함
	channel, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
	}
	defer channel.Close()

	// (4) protoc가 생성한 클라이언트 스텁을 생성함
	stub := pb.NewMyServiceClient(channel)

	// (5) 요청 메시지를 생성하고 값 할당
	request := &pb.MyNumber{
		Value: 4,
	}

	// (6) 원격 함수를 호출함
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := stub.MyFunction(ctx, request)
	if err != nil {
		fmt.Printf("Error when calling MyFunction: %v", err)
	}

	// (7) 결과를 활용함
	fmt.Printf("gRPC result: %d\n", response.Value)
}
