#!/bin/bash

echo "generate server"
protoc \
    --go_out=server \
    --go_opt=paths=source_relative \
    --go-grpc_out=server \
    --go-grpc_opt=paths=source_relative \
    proto/greeter.proto

echo "generate go-clinet"
protoc \
    --go_out=client/go \
    --go_opt=paths=source_relative \
    --go-grpc_out=client/go \
    --go-grpc_opt=paths=source_relative \
    proto/greeter.proto

echo "generate flutter-client"
protoc \
    --dart_out=grpc:client/flutter/lib/proto \
    -Iproto proto/greeter.proto