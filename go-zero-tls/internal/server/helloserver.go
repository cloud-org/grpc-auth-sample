// Code generated by goctl. DO NOT EDIT!
// Source: hello.proto

package server

import (
	"context"

	"grpc-sample/go-zero-tls/hello"
	"grpc-sample/go-zero-tls/internal/logic"
	"grpc-sample/go-zero-tls/internal/svc"
)

type HelloServer struct {
	svcCtx *svc.ServiceContext
}

func NewHelloServer(svcCtx *svc.ServiceContext) *HelloServer {
	return &HelloServer{
		svcCtx: svcCtx,
	}
}

func (s *HelloServer) Ping(ctx context.Context, in *hello.Request) (*hello.Response, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}
