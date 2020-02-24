package test_clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	grpcclients "github.com/pip-services3-go/pip-services3-grpc-go/clients"
	testgrpc "github.com/pip-services3-go/pip-services3-grpc-go/test"
	rpcclients "github.com/pip-services3-go/pip-services3-rpc-go/clients"
)

type DummyCommandableGrpcClient struct {
	*grpcclients.CommandableGrpcClient
}

func NewDummyCommandableGrpcClient() *DummyCommandableGrpcClient {
	dcgc := DummyCommandableGrpcClient{}
	dcgc.CommandableGrpcClient = grpcclients.NewCommandableGrpcClient("dummy")
	return &dcgc
}

func (c *DummyCommandableGrpcClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *testgrpc.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	calValue, calErr := c.CallCommand("get_dummies", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}

	convRes, err := rpcclients.ConvertComandResult(calValue, reflect.TypeOf(&testgrpc.DummyDataPage{}))
	result, _ = convRes.(*testgrpc.DummyDataPage)
	return result, err
}

func (c *DummyCommandableGrpcClient) GetDummyById(correlationId string, dummyId string) (result *testgrpc.Dummy, err error) {

	params := cdata.NewEmptyStringValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand("get_dummy_by_id", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	convRes, err := rpcclients.ConvertComandResult(calValue, reflect.TypeOf(&testgrpc.Dummy{}))
	result, _ = convRes.(*testgrpc.Dummy)
	return result, err
}

func (c *DummyCommandableGrpcClient) CreateDummy(correlationId string, dummy testgrpc.Dummy) (result *testgrpc.Dummy, err error) {

	bodyMap := make(map[string]interface{})
	bodyMap["dummy"] = dummy
	calValue, calErr := c.CallCommand("create_dummy", correlationId, bodyMap)
	if calErr != nil {
		return nil, calErr
	}

	convRes, err := rpcclients.ConvertComandResult(calValue, reflect.TypeOf(&testgrpc.Dummy{}))
	result, _ = convRes.(*testgrpc.Dummy)
	return result, err
}

func (c *DummyCommandableGrpcClient) UpdateDummy(correlationId string, dummy testgrpc.Dummy) (result *testgrpc.Dummy, err error) {

	bodyMap := make(map[string]interface{})
	bodyMap["dummy"] = dummy
	calValue, calErr := c.CallCommand("update_dummy", correlationId, bodyMap)
	if calErr != nil {
		return nil, calErr
	}
	convRes, err := rpcclients.ConvertComandResult(calValue, reflect.TypeOf(&testgrpc.Dummy{}))
	result, _ = convRes.(*testgrpc.Dummy)
	return result, err
}

func (c *DummyCommandableGrpcClient) DeleteDummy(correlationId string, dummyId string) (result *testgrpc.Dummy, err error) {

	params := cdata.NewEmptyStringValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand("delete_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	convRes, err := rpcclients.ConvertComandResult(calValue, reflect.TypeOf(&testgrpc.Dummy{}))
	result, _ = convRes.(*testgrpc.Dummy)
	return result, err
}
