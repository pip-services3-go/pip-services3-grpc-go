package build

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
	grpcservices "github.com/pip-services3-go/pip-services3-grpc-go/services"
)

/*
Creates GRPC components by their descriptors.
See Factory
See GrpcEndpoint
See HeartbeatGrpcService
See StatusGrpcService
*/

type DefaultGrpcFactory struct {
	*cbuild.Factory
	Descriptor             *cref.Descriptor
	GrpcEndpointDescriptor *cref.Descriptor
	//  StatusServiceDescriptor *cref.Descriptor = new Descriptor("pip-services", "status-service", "grpc", "*", "1.0");
	//  HeartbeatServiceDescriptor *cref.Descriptor = new Descriptor("pip-services", "heartbeat-service", "grpc", "*", "1.0");
}

/*
	Create a new instance of the factory.
*/
func NewDefaultGrpcFactory() *DefaultGrpcFactory {

	dgf := DefaultGrpcFactory{}
	dgf.Descriptor = cref.NewDescriptor("pip-services", "factory", "grpc", "default", "1.0")
	dgf.GrpcEndpointDescriptor = cref.NewDescriptor("pip-services", "endpoint", "grpc", "*", "1.0")
	//  dgf.StatusServiceDescriptor  = cref.NewDescriptor("pip-services", "status-service", "grpc", "*", "1.0");
	//  dgf.HeartbeatServiceDescriptor = cref.NewDescriptor("pip-services", "heartbeat-service", "grpc", "*", "1.0");

	dgf.RegisterType(dgf.GrpcEndpointDescriptor, grpcservices.NewGrpcEndpoint)
	// dgf.RegisterType(dgf.HeartbeatServiceDescriptor, grpcservices.NewHeartbeatGrpcService);
	// dgf.RegisterType(dgf.StatusServiceDescriptor, grpcservices.NewStatusGrpcService);
	return &dgf
}
