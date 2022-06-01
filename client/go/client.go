package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"time"

	pb "deniffel.com/grpc_nginx_https/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "api.example.com:443"
)

func main() {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithBlock())
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
