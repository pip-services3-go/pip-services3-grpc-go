package test_services

import (
	"context"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	grpcservices "github.com/pip-services3-go/pip-services3-grpc-go/services"
	testgrpc "github.com/pip-services3-go/pip-services3-grpc-go/test"
	"github.com/pip-services3-go/pip-services3-grpc-go/test/protos"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func TestDummyGrpcServiceConnection(t *testing.T) {

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	var Dummy1 testgrpc.Dummy
	var Dummy2 testgrpc.Dummy

	var service *DummyGrpcService
	var endpoint *grpcservices.GrpcEndpoint
	var client protos.DummiesClient
	ctrl := testgrpc.NewDummyController()

	endpoint = grpcservices.NewGrpcEndpoint()
	endpoint.Configure(grpcConfig)

	service = NewDummyGrpcService()
	service.Configure(grpcConfig)

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "endpoint", "grpc", "default", "1.0"), endpoint,
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "grpc", "default", "1.0"), service,
	)
	service.SetReferences(references)

	eErr := endpoint.Open("")
	if eErr != nil {
		grpclog.Fatalf("Fail to open endpoint: %v", eErr)
	}

	service.Open("")

	defer service.Close("")

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:3000", opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()
	client = protos.NewDummiesClient(conn)

	assert.True(t, endpoint.IsOpen())

	Dummy1 = testgrpc.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	Dummy2 = testgrpc.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	// Test CRUD Operations
	// Create first dummy
	protoDummy := protos.Dummy{}
	protoDummy.Id = Dummy1.Id
	protoDummy.Key = Dummy1.Key
	protoDummy.Content = Dummy1.Content
	request := protos.DummyObjectRequest{Dummy: &protoDummy}
	dummy, err := client.CreateDummy(context.TODO(), &request)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, protoDummy.Content, dummy.Content)
	assert.Equal(t, protoDummy.Key, dummy.Key)

	dummy1 := dummy

	// Create another dummy
	protoDummy.Id = Dummy2.Id
	protoDummy.Key = Dummy2.Key
	protoDummy.Content = Dummy2.Content
	request = protos.DummyObjectRequest{Dummy: &protoDummy}
	dummy, err = client.CreateDummy(context.TODO(), &request)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, protoDummy.Content, dummy.Content)
	assert.Equal(t, protoDummy.Key, dummy.Key)

	// Get all dummies
	requestPage := protos.DummiesPageRequest{}
	dummies, err := client.GetDummies(context.TODO(), &requestPage)
	assert.Nil(t, err)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	protoDummy.Id = dummy1.Id
	protoDummy.Key = dummy1.Key
	protoDummy.Content = dummy1.Content

	request = protos.DummyObjectRequest{Dummy: &protoDummy}
	dummy, err = client.UpdateDummy(context.TODO(), &request)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, dummy1.Key)

	// Delete dummy
	idRequest := protos.DummyIdRequest{DummyId: dummy1.Id}
	dummy, err = client.DeleteDummyById(context.TODO(), &idRequest)
	assert.Nil(t, err)

	// Try to get delete dummy
	idRequest = protos.DummyIdRequest{DummyId: dummy1.Id}
	dummy, err = client.GetDummyById(context.TODO(), &idRequest)
	assert.Nil(t, err)
}
