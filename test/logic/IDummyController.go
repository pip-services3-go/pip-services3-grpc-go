package test_logic

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	tdata "github.com/pip-services3-go/pip-services3-grpc-go/test/data"
)

type IDummyController interface {
	GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *tdata.DummyDataPage, err error)
	GetOneById(correlationId string, id string) (result *tdata.Dummy, err error)
	Create(correlationId string, entity tdata.Dummy) (result *tdata.Dummy, err error)
	Update(correlationId string, entity tdata.Dummy) (result *tdata.Dummy, err error)
	DeleteById(correlationId string, id string) (result *tdata.Dummy, err error)
}
