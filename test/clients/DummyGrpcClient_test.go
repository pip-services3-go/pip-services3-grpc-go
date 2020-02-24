package test_clients

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	testgrpc "github.com/pip-services3-go/pip-services3-grpc-go/test"
	testservices "github.com/pip-services3-go/pip-services3-grpc-go/test/services"
)

func TestDummyGrpcClient(t *testing.T) {

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	var service *testservices.DummyGrpcService
	var client *DummyGrpcClient
	var fixture *DummyClientFixture

	ctrl := testgrpc.NewDummyController()

	service = testservices.NewDummyGrpcService()
	service.Configure(grpcConfig)

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "grpc", "default", "1.0"), service,
	)
	service.SetReferences(references)

	service.Open("")

	defer service.Close("")

	client = NewDummyGrpcClient()
	fixture = NewDummyClientFixture(client)

	client.Configure(grpcConfig)
	client.SetReferences(cref.NewEmptyReferences())
	client.Open("")
	defer client.Close("")

	t.Run("CRUD Operations", fixture.TestCrudOperations)

}
