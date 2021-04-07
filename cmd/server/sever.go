package main

import (
	"fmt"
	"net"

	pb "grpc-sample/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials" // 引入 grpc 认证包
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC 服务地址
	Address = "127.0.0.1:50052"
)

// 定义 helloService 并实现约定的接口
type helloService struct{}

// HelloService Hello 服务
var HelloService = helloService{}

// SayHello 实现 Hello 服务接口
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.", in.Name)

	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// TLS 认证
	creds, err := credentials.NewServerTLSFromFile("../../tls/server.pem", "../../tls/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	// 实例化 grpc Server, 并开启 TLS 认证
	s := grpc.NewServer(grpc.Creds(creds))

	// 注册 HelloService
	pb.RegisterHelloServer(s, HelloService)

	grpclog.Println("Listen on " + Address + " with TLS")

	s.Serve(listen)
}
