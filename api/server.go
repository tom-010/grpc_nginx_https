package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "deniffel.com/grpc_nginx_https/proto"
	"google.golang.org/grpc"
)

const (
	port = "0.0.0.0:50051"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name + " !"}, nil
}

func (s *GreeterServer) SayRepeatHello(in *pb.RepeatHelloRequest, stream pb.Greeter_SayRepeatHelloServer) error {
	for i := 1; i <= int(in.Count); i++ {
		stream.Send(&pb.HelloReply{Message: fmt.Sprintf("Hello %d!", i)})
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listn: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &GreeterServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
