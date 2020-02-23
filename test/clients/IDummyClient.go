package test_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	testgrpc "github.com/pip-services3-go/pip-services3-grpc-go/test"
)

type IDummyClient interface {
	GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *testgrpc.DummyDataPage, err error)
	GetDummyById(correlationId string, dummyId string) (result *testgrpc.Dummy, err error)
	CreateDummy(correlationId string, dummy testgrpc.Dummy) (result *testgrpc.Dummy, err error)
	UpdateDummy(correlationId string, dummy testgrpc.Dummy) (result *testgrpc.Dummy, err error)
	DeleteDummy(correlationId string, dummyId string) (result *testgrpc.Dummy, err error)
}
