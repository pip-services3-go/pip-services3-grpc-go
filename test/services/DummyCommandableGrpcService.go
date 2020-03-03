package test_services

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	grpcservices "github.com/pip-services3-go/pip-services3-grpc-go/services"
)

type DummyCommandableGrpcService struct {
	*grpcservices.CommandableGrpcService
}

func NewDummyCommandableGrpcService() *DummyCommandableGrpcService {

	c := DummyCommandableGrpcService{}
	c.CommandableGrpcService = grpcservices.NewCommandableGrpcService("dummy")
	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &c
}
