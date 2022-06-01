package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "deniffel.com/grpc_nginx_https/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	port = "0.0.0.0:50051"
)

var (
	username           = os.Getenv("USERNAME")
	password           = os.Getenv("PASSWORD")
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid credentials")
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	token, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(md["authorization"][0], "Basic "))
	if err != nil {
		return nil, errMissingMetadata
	}
	username, password, _ := strings.Cut(string(token), ":")
	log.Println(password)
	return &pb.HelloReply{Message: "Hello " + in.Name + " ! current user:" + username}, nil
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
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidBasicCredentials),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(s, &GreeterServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")
	log.Println(token)
	return token == base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func ensureValidBasicCredentials(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}
