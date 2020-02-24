package test_clients

import (
	"encoding/json"

	"github.com/pip-services3-go/pip-services3-commons-go/convert"
	"github.com/pip-services3-go/pip-services3-commons-go/errors"
	testgrpc "github.com/pip-services3-go/pip-services3-grpc-go/test"
	testproto "github.com/pip-services3-go/pip-services3-grpc-go/test/protos"
)

func fromError(err error) *testproto.ErrorDescription {
	if err == nil {
		return nil
	}

	desc := errors.ErrorDescriptionFactory.Create(err)
	obj := &testproto.ErrorDescription{
		//Type:          desc.Type,
		Category:      desc.Category,
		Code:          desc.Code,
		CorrelationId: desc.CorrelationId,
		Status:        convert.StringConverter.ToString(desc.Status),
		Message:       desc.Message,
		Cause:         desc.Cause,
		StackTrace:    desc.StackTrace,
		Details:       fromMap(desc.Details),
	}

	return obj
}

func toError(obj *testproto.ErrorDescription) error {
	if obj == nil || (obj.Category == "" && obj.Message == "") {
		return nil
	}

	description := &errors.ErrorDescription{
		//Type:          obj.Type,
		Category:      obj.Category,
		Code:          obj.Code,
		CorrelationId: obj.CorrelationId,
		Status:        convert.IntegerConverter.ToInteger(obj.Status),
		Message:       obj.Message,
		Cause:         obj.Cause,
		StackTrace:    obj.StackTrace,
		Details:       toMap(obj.Details),
	}

	return errors.ApplicationErrorFactory.Create(description)
}

func fromMap(val map[string]interface{}) map[string]string {
	r := map[string]string{}

	for k, v := range val {
		r[k] = convert.ToString(v)
	}

	return r
}

func toMap(val map[string]string) map[string]interface{} {
	var r map[string]interface{}

	for k, v := range val {
		r[k] = v
	}

	return r
}

func toJson(value interface{}) string {
	if value == nil {
		return ""
	}

	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(b[:])
}

func fromJson(value string) interface{} {
	if value == "" {
		return nil
	}

	var m interface{}
	json.Unmarshal([]byte(value), &m)
	return m
}

func fromDummy(in *testgrpc.Dummy) *testproto.Dummy {
	if in == nil {
		return nil
	}

	obj := &testproto.Dummy{
		Id:      in.Id,
		Key:     in.Key,
		Content: in.Content,
	}

	return obj
}

func toDummy(obj *testproto.Dummy) *testgrpc.Dummy {
	if obj == nil {
		return nil
	}

	dummy := &testgrpc.Dummy{
		Id:      obj.Id,
		Key:     obj.Key,
		Content: obj.Content,
	}

	return dummy
}

func toDummiesPage(obj *testproto.DummiesPage) *testgrpc.DummyDataPage {
	if obj == nil {
		return nil
	}

	dummies := make([]testgrpc.Dummy, len(obj.Data))
	for i, v := range obj.Data {
		dummies[i] = *toDummy(v)
	}
	page := testgrpc.NewDummyDataPage(&obj.Total, dummies)

	return page
}
