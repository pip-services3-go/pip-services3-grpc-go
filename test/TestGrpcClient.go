package test

import (
	"github.com/pip-services3-go/pip-services3-grpc-go/clients"
)

type TestGrpcClient struct {
	clients.GrpcClient
}

func NewTestRestClient(name string) *TestGrpcClient {
	c := &TestGrpcClient{}
	c.GrpcClient = *clients.NewGrpcClient(name)
	return c
}
