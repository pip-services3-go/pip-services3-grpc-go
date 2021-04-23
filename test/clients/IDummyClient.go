package test_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	tdata "github.com/pip-services3-go/pip-services3-grpc-go/test/data"
)

type IDummyClient interface {
	GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error)
	GetDummyById(correlationId string, dummyId string) (result *tdata.Dummy, err error)
	CreateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error)
	UpdateDummy(correlationId string, dummy tdata.Dummy) (result *tdata.Dummy, err error)
	DeleteDummy(correlationId string, dummyId string) (result *tdata.Dummy, err error)
}
