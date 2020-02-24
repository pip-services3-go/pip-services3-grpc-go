package test_services

import (
	"context"
	"encoding/json"

	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	grpcservices "github.com/pip-services3-go/pip-services3-grpc-go/services"
	grpctest "github.com/pip-services3-go/pip-services3-grpc-go/test"
	"github.com/pip-services3-go/pip-services3-grpc-go/test/protos"
	"google.golang.org/grpc"
)

type DummyGrpcService struct {
	*grpcservices.GrpcService
	controller    grpctest.IDummyController
	numberOfCalls int64
}

func NewDummyGrpcService() *DummyGrpcService {
	dgs := DummyGrpcService{}
	dgs.GrpcService = grpcservices.NewGrpcService("dummies.Dummies")
	dgs.GrpcService.IRegisterable = &dgs
	dgs.numberOfCalls = 0
	dgs.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &dgs
}

func (c *DummyGrpcService) SetReferences(references cref.IReferences) {
	c.GrpcService.SetReferences(references)
	resolv, err := c.DependencyResolver.GetOneRequired("controller")
	if err == nil && resolv != nil {
		c.controller = resolv.(grpctest.IDummyController)
		return
	}
	panic("Can't resolve 'controller' reference")
}

func (c *DummyGrpcService) GetNumberOfCalls() int64 {
	return c.numberOfCalls
}

func (c *DummyGrpcService) incrementNumberOfCalls(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	m, err := handler(ctx, req)
	if err != nil {
		c.Logger.Error("", err, "RPC failed with error %v", err.Error())
	}
	c.numberOfCalls++
	return m, err
}

func (c *DummyGrpcService) Open(correlationId string) error {
	// Add interceptors
	c.Endpoint.AddInterceptors(grpc.UnaryInterceptor(c.incrementNumberOfCalls))

	return c.GrpcService.Open(correlationId)
}

func (c *DummyGrpcService) GetDummies(ctx context.Context, req *protos.DummiesPageRequest) (*protos.DummiesPage, error) {

	// Schema := cvalid.NewObjectSchema().
	// 	WithOptionalProperty("Paging", cvalid.NewPagingParamsSchema()).
	// 	WithOptionalProperty("Filter", cvalid.NewFilterParamsSchema())

	// validateErr := Schema.ValidateAndReturnError("", *req, false)

	// if validateErr != nil {
	// 	return nil, validateErr
	// }

	filter := cdata.NewFilterParamsFromValue(req.GetFilter())
	paging := cdata.NewEmptyPagingParams()
	if req.Paging != nil {
		paging = cdata.NewPagingParams(req.Paging.GetSkip(), req.Paging.GetTake(), req.Paging.GetTotal())
	}
	data, err := c.controller.GetPageByFilter(
		req.CorrelationId,
		filter,
		paging,
	)
	if err != nil || data == nil {
		return nil, err
	}

	result := protos.DummiesPage{}
	result.Total = *data.Total
	for _, v := range data.Data {
		buf := protos.Dummy{}
		bytes, _ := json.Marshal(v)
		json.Unmarshal(bytes, &buf)
		result.Data = append(result.Data, &buf)
	}

	return &result, err
}

func (c *DummyGrpcService) GetDummyById(ctx context.Context, req *protos.DummyIdRequest) (*protos.Dummy, error) {

	Schema := cvalid.NewObjectSchema().
		WithRequiredProperty("DummyId", cconv.String)

	validateErr := Schema.ValidateAndReturnError("", *req, false)

	if validateErr != nil {
		return nil, validateErr
	}

	data, err := c.controller.GetOneById(
		req.CorrelationId,
		req.DummyId,
	)
	if err != nil {
		return nil, err
	}
	result := protos.Dummy{}
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &result)

	return &result, nil

}

func (c *DummyGrpcService) CreateDummy(ctx context.Context, req *protos.DummyObjectRequest) (*protos.Dummy, error) {

	Schema := cvalid.NewObjectSchema().
		WithRequiredProperty("Dummy", grpctest.NewDummySchema())

	validateErr := Schema.ValidateAndReturnError("", *req, false)

	if validateErr != nil {
		return nil, validateErr
	}

	dummy := grpctest.Dummy{}
	bytes, _ := json.Marshal(req.Dummy)
	json.Unmarshal(bytes, &dummy)

	data, err := c.controller.Create(
		req.CorrelationId,
		dummy,
	)

	if err != nil || data == nil {
		return nil, err
	}
	result := protos.Dummy{}
	bytes, _ = json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil

}

func (c *DummyGrpcService) UpdateDummy(ctx context.Context, req *protos.DummyObjectRequest) (*protos.Dummy, error) {

	Schema := cvalid.NewObjectSchema().
		WithRequiredProperty("Dummy", grpctest.NewDummySchema())

	validateErr := Schema.ValidateAndReturnError("", *req, false)

	if validateErr != nil {
		return nil, validateErr
	}

	dummy := grpctest.Dummy{}
	bytes, _ := json.Marshal(req.Dummy)
	json.Unmarshal(bytes, &dummy)

	data, err := c.controller.Update(
		req.CorrelationId,
		dummy,
	)

	if err != nil || data == nil {
		return nil, err
	}
	result := protos.Dummy{}
	bytes, _ = json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil
}

func (c *DummyGrpcService) DeleteDummyById(ctx context.Context, req *protos.DummyIdRequest) (*protos.Dummy, error) {

	Schema := cvalid.NewObjectSchema().
		WithRequiredProperty("DummyId", cconv.String)

	validateErr := Schema.ValidateAndReturnError("", *req, false)

	if validateErr != nil {
		return nil, validateErr
	}

	data, err := c.controller.DeleteById(
		req.CorrelationId,
		req.DummyId,
	)
	if err != nil || data == nil {
		return nil, err
	}
	result := protos.Dummy{}
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil
}

func (c *DummyGrpcService) Register() {

	protos.RegisterDummiesServer(c.Endpoint.GetServer(), c)

}
