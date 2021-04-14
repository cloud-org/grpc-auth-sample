package logic

import (
	"context"
	"fmt"

	"grpc-sample/go-zero-tls/hello"
	"grpc-sample/go-zero-tls/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *hello.Request) (*hello.Response, error) {

	return &hello.Response{
		Pong: fmt.Sprintf("hello: %s", in.Ping),
	}, nil
}
