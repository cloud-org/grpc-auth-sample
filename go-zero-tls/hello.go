package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"flag"
	"fmt"
	"log"

	"grpc-sample/go-zero-tls/hello"
	"grpc-sample/go-zero-tls/internal/config"
	"grpc-sample/go-zero-tls/internal/server"
	"grpc-sample/go-zero-tls/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	configFile = flag.String("f", "etc/hello.yaml", "the config file")

	//go:embed tls/server.pem
	serverPem []byte

	//go:embed tls/server.key
	serverKey []byte

	//go:embed tls/ca.pem
	caPem []byte
)

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewHelloServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		hello.RegisterHelloServer(grpcServer, srv)
	})

	// add creds
	cred := getCreds()
	s.AddOptions(grpc.Creds(cred))

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

//getCreds 添加凭证
func getCreds() credentials.TransportCredentials {
	// cert, err := tls.LoadX509KeyPair("../../tls/server.pem", "../../tls/server.key")

	cert, err := tls.X509KeyPair(serverPem, serverKey)
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	// ca, err := os.ReadFile("../../tls/ca.pem")
	// if err != nil {
	// 	log.Fatalf("ioutil.ReadFile err: %v", err)
	// }

	if ok := certPool.AppendCertsFromPEM(caPem); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return creds
}
