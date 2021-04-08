package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"grpc-sample/keys"
	pb "grpc-sample/proto" // 引入 proto 包
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials" // 引入 grpc 认证包
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC 服务地址
	Address = "127.0.0.1:50052"
)

type Auth struct {
	AppKey    string `json:"app_key"`
	SecretKey string `json:"secret_key"`
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "secret_key": a.SecretKey}, nil
}

func (a *Auth) RequireTransportSecurity() bool { // 是否 TLS 认证
	// 如果设置了 true 但是 client 连接的时候使用 grpc.WithInsecure() 会直接 Fatal
	// grpc: the credentials require transport level security (use grpc.WithTransportCredentials() to set)
	return true
}

func main() {
	// TLS 连接
	cert, err := tls.LoadX509KeyPair("../../tls/client.pem", "../../tls/client.key")
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
		ServerName:   "localhost",
		RootCAs:      certPool,
	})

	auth := Auth{
		AppKey:    keys.GenAppKey(),
		SecretKey: keys.GenSecretKey(),
	}

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&auth))
	//conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用方法
	req := &pb.HelloRequest{Name: "ashing"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		grpclog.Fatalln(err)
	}

	grpclog.Println(res.Message)
}
