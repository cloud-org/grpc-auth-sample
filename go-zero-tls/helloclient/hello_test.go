package helloclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"testing"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestPing(t *testing.T) {

	// 添加证书设置
	cred := getCreds()
	client := zrpc.MustNewClient(
		zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"127.0.0.1:2379"},
				Key:   "hello.rpc",
			},
		},
		zrpc.WithDialOption(grpc.WithTransportCredentials(cred)),
	)

	h := NewHello(client)
	resp, err := h.Ping(context.TODO(), &Request{
		Ping: "ashing",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp.GetPong())
}

//getCreds 添加凭证
func getCreds() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("../tls/client.pem", "../tls/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("../tls/ca.pem")
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

	return creds
}
