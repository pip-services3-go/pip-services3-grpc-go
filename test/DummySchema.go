package test_rpc

// import { TypeCode } from 'pip-services3-commons-node';
// import { ObjectSchema } from 'pip-services3-commons-node';
import (
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

type DummySchema struct {
	cvalid.ObjectSchema
}

func NewDummySchema() *DummySchema {
	ds := DummySchema{}
	ds.WithOptionalProperty("id", cconv.String)
	ds.WithRequiredProperty("key", cconv.String)
	ds.WithOptionalProperty("content", cconv.String)
	return &ds
}
