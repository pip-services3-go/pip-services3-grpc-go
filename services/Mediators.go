package services

import (
	"context"

	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"
)

// InvokeComandMediator Helper class for implements invoke method in CommandableGrpc
type InvokeComandMediator struct {
	InvokeFunc func(ctx context.Context, request *grpcproto.InvokeRequest) (response *grpcproto.InvokeReply, err error)
}

func (c *InvokeComandMediator) Invoke(ctx context.Context, request *grpcproto.InvokeRequest) (response *grpcproto.InvokeReply, err error) {
	return c.InvokeFunc(ctx, request)
}
