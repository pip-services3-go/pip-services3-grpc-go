package test_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	grpcclients "github.com/pip-services3-go/pip-services3-grpc-go/clients"
	tdata "github.com/pip-services3-go/pip-services3-grpc-go/test/data"
	testproto "github.com/pip-services3-go/pip-services3-grpc-go/test/protos"
)

type DummyGrpcClient struct {
	*grpcclients.GrpcClient
}

func NewDummyGrpcClient() *DummyGrpcClient {
	dgc := DummyGrpcClient{}
	dgc.GrpcClient = grpcclients.NewGrpcClient("dummies.Dummies")
	return &dgc
}

func (c *DummyGrpcClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error) {

	req := &testproto.DummiesPageRequest{
		CorrelationId: correlationId,
	}
	if filter != nil {
		req.Filter = filter.Value()
	}
	if paging != nil {
		req.Paging = &testproto.PagingParams{
			Skip:  paging.GetSkip(0),
			Take:  (int32)(paging.GetTake(100)),
			Total: paging.Total,
		}
	}
	reply := new(testproto.DummiesPage)
	err = c.Call("get_dummies", correlationId, req, reply)
	c.Instrument(correlationId, "dummy.get_page_by_filter")
	if err != nil {
		return nil, err
	}
	result = toDummiesPage(reply)
	return result, nil
}

func (c *DummyGrpcClient) GetDummyById(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	req := &testproto.DummyIdRequest{
		CorrelationId: correlationId,
		DummyId:       dummyId,
	}

	reply := new(testproto.Dummy)
	err = c.Call("get_dummy_by_id", correlationId, req, reply)
	c.Instrument(correlationId, "dummy.get_one_by_id")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}

func (c *DummyGrpcClient) CreateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {

	req := &testproto.DummyObjectRequest{
		CorrelationId: correlationId,
		Dummy:         fromDummy(&dummy),
	}

	reply := new(testproto.Dummy)
	err = c.Call("create_dummy", correlationId, req, reply)
	c.Instrument(correlationId, "dummy.create")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}

func (c *DummyGrpcClient) UpdateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error) {
	req := &testproto.DummyObjectRequest{
		CorrelationId: correlationId,
		Dummy:         fromDummy(&dummy),
	}
	reply := new(testproto.Dummy)
	err = c.Call("update_dummy", correlationId, req, reply)
	c.Instrument(correlationId, "dummy.update")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}

func (c *DummyGrpcClient) DeleteDummy(correlationId string, dummyId string) (result *tdata.Dummy, err error) {

	req := &testproto.DummyIdRequest{
		CorrelationId: correlationId,
		DummyId:       dummyId,
	}

	reply := new(testproto.Dummy)
	c.Call("delete_dummy_by_id", correlationId, req, reply)
	c.Instrument(correlationId, "dummy.delete_by_id")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}
