package test_services

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	grpcservices "github.com/pip-services3-go/pip-services3-grpc-go/services"
	"github.com/stretchr/testify/assert"
)

func TestGrpcEndpoint(t *testing.T) {

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", 3000,
	)

	var endpoint *grpcservices.GrpcEndpoint

	endpoint = grpcservices.NewGrpcEndpoint()
	endpoint.Configure(grpcConfig)

	endpoint.Open("")
	assert.True(t, endpoint.IsOpen())
	endpoint.Close("")
}
