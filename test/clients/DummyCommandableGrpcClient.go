package test_clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	grpcclients "github.com/pip-services3-go/pip-services3-grpc-go/clients"
	tdata "github.com/pip-services3-go/pip-services3-grpc-go/test/data"
)

var (
	dummyDataPageType = reflect.TypeOf(&tdata.DummyDataPage{})
	dummyType         = reflect.TypeOf(&tdata.Dummy{})
)

type DummyCommandableGrpcClient struct {
	*grpcclients.CommandableGrpcClient
}

func NewDummyCommandableGrpcClient() *DummyCommandableGrpcClient {
	dcgc := DummyCommandableGrpcClient{}
	dcgc.CommandableGrpcClient = grpcclients.NewCommandableGrpcClient("dummy")
	return &dcgc
}

func (c *DummyCommandableGrpcClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	calValue, calErr := c.CallCommand(dummyDataPageType, "get_dummies", correlationId, cdata.NewAnyValueMapFromValue(params.Value()))
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.DummyDataPage)
	return result, err
}

func (c *DummyCommandableGrpcClient) GetDummyById(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "get_dummy_by_id", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableGrpcClient) CreateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "create_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableGrpcClient) UpdateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "update_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}

func (c *DummyCommandableGrpcClient) DeleteDummy(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "delete_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*tdata.Dummy)
	return result, err
}
