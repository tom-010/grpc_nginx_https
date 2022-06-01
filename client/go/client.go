package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "deniffel.com/grpc_nginx_https/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Tom"})
	if err != nil {
		log.Fatalf("SayHello error: %v", err)
	}
	log.Printf(r.Message)

	greetingStream, _ := c.SayRepeatHello(ctx, &pb.RepeatHelloRequest{Name: "Tom", Count: 5})

	for {
		r, err = greetingStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error in stream: %v", err)
		}
		log.Println(r.Message)
	}

}
