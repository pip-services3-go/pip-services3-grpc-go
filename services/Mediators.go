package services

import (
	"context"

	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"
)

type InvokeComandMediator struct {
	InvokeFunc func(ctx context.Context, request *grpcproto.InvokeRequest) (response *grpcproto.InvokeReply, err error)
}

func (c *InvokeComandMediator) Invoke(ctx context.Context, request *grpcproto.InvokeRequest) (response *grpcproto.InvokeReply, err error) {
	return c.InvokeFunc(ctx, request)
}
