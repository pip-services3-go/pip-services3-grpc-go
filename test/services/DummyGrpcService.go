package test_services

import (
	"context"
	"encoding/json"

	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	grpcservices "github.com/pip-services3-go/pip-services3-grpc-go/services"
	tdata "github.com/pip-services3-go/pip-services3-grpc-go/test/data"
	tlogic "github.com/pip-services3-go/pip-services3-grpc-go/test/logic"
	"github.com/pip-services3-go/pip-services3-grpc-go/test/protos"
	"google.golang.org/grpc"
)

type DummyGrpcService struct {
	grpcservices.GrpcService
	controller    tlogic.IDummyController
	numberOfCalls int64
}

func NewDummyGrpcService() *DummyGrpcService {
	c := &DummyGrpcService{}
	c.GrpcService = *grpcservices.InheritGrpcService(c, "dummies.Dummies")
	c.numberOfCalls = 0
	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyGrpcService) SetReferences(references cref.IReferences) {
	c.GrpcService.SetReferences(references)
	resolv, err := c.DependencyResolver.GetOneRequired("controller")
	if err == nil && resolv != nil {
		c.controller = resolv.(tlogic.IDummyController)
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
	c.RegisterUnaryInterceptor(c.incrementNumberOfCalls)
	return c.GrpcService.Open(correlationId)
}

func (c *DummyGrpcService) GetDummies(ctx context.Context, req *protos.DummiesPageRequest) (*protos.DummiesPage, error) {

	validateErr := c.ValidateRequest(req,
		&cvalid.NewObjectSchema().
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()).
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

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

	// validation
	validateErr := c.ValidateRequest(req,
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema)

	if validateErr != nil {
		return nil, validateErr
	}
	// ==================================

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

	// validation
	validateErr := c.ValidateRequest(req,
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", tdata.NewDummySchema()).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

	dummy := tdata.Dummy{}
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

	validateErr := c.ValidateRequest(req,
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", tdata.NewDummySchema()).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

	dummy := tdata.Dummy{}
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

	validateErr := c.ValidateRequest(req,
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema)

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
