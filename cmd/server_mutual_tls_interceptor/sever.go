package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"

	pb "grpc-sample/proto"

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

func Check(ctx context.Context) error {
	//从上下文中获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	grpclog.Infof("md is %+v\n", md)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "获取 Token 失败")
	}
	var (
		appKey    string
		secretKey string
	)
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["secret_key"]; ok {
		secretKey = value[0]
	}

	if len(appKey) != 16 || len(secretKey) != 32 {
		return status.Errorf(codes.Unauthenticated, "Token err")
	}

	grpclog.Infof("auth is %v, %v\n", appKey, secretKey)
	return nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	interceptor := getInterceptor()
	creds := getCreds(err)
	// 实例化 grpc Server, 并开启 TLS 认证
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(interceptor))

	// 注册 HelloService
	pb.RegisterHelloServer(s, HelloService)

	grpclog.Println("Listen on " + Address + " with TLS")

	s.Serve(listen)
}

//getInterceptor 添加拦截器
func getInterceptor() grpc.UnaryServerInterceptor {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//拦截普通方法请求，验证 Token
		//req proto.HelloRequest
		//info {"Server": "main.helloService", "FullMethod": "/hello.Hello/SayHello"}
		//handler real rpc func
		err = Check(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	return interceptor
}

//getCreds 添加凭证
func getCreds(err error) credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("../../tls/server.pem", "../../tls/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("../../tls/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return creds
}
