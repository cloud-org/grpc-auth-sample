package main

import (
	pb "grpc-sample/proto" // 引入 proto 包

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials" // 引入 grpc 认证包
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC 服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	// TLS 连接
	creds, err := credentials.NewClientTLSFromFile("../../tls/server.pem", "localhost")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用方法
	req := &pb.HelloRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		grpclog.Fatalln(err)
	}

	grpclog.Println(res.Message)
}
