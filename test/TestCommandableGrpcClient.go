package test

import (
	"github.com/pip-services3-go/pip-services3-grpc-go/clients"
)

type TestCommandableGrpcClient struct {
	clients.CommandableGrpcClient
}

func NewTestCommandableGrpcClient(name string) *TestCommandableGrpcClient {
	c := &TestCommandableGrpcClient{}
	c.CommandableGrpcClient = *clients.NewCommandableGrpcClient(name)
	return c
}
