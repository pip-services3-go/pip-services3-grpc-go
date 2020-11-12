package clients

import (
	"encoding/json"
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"
	rpcclients "github.com/pip-services3-go/pip-services3-rpc-go/clients"
)

/*
CommandableGrpcClient abstract client that calls commandable GRPC service.

Commandable services are generated automatically for ICommandable objects.
Each command is exposed as Invoke method that receives all parameters as args.

Configuration parameters:

- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- options:
  - retries:               number of retries (default: 3)
  - connect_timeout:       connection timeout in milliseconds (default: 10 sec)
  - timeout:               invocation timeout in milliseconds (default: 10 sec)

 References:

- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
- *:counters:*:*:1.0         (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0        (optional) IDiscovery services to resolve connection

Example:

    type MyCommandableGrpcClient struct {
	 *CommandableGrpcClient
       ...
	}
        func (c * MyCommandableGrpcClient) GetData(correlationId string, id string) (result *MyData, err error) {
           params := cdata.NewEmptyStringValueMap()
			params.Put("id", id)

			calValue, calErr := c.CallCommand(MyDataType, "get_mydata_by_id", correlationId, params)
			if calErr != nil {
				return nil, calErr
			}
			result, _ = calValue.(*testgrpc.MyData)
			return result, err
        }
        ...

    client := NewMyCommandableGrpcClient();
    client.Configure(cconf.NewConfigParamsFromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080,
    ));

	result, err := client.GetData("123", "1")
    ...
*/
type CommandableGrpcClient struct {
	*GrpcClient
	//The service name
	Name string
}

// NewCommandableGrpcClient method are creates a new instance of the client.
// Parameters:
// - name     a service name.
func NewCommandableGrpcClient(name string) *CommandableGrpcClient {
	c := CommandableGrpcClient{}
	c.GrpcClient = NewGrpcClient("commandable.Commandable")
	c.Name = name
	return &c
}

// CallCommand method are calls a remote method via GRPC commadable protocol.
// The call is made via Invoke method and all parameters are sent in args object.
// The complete route to remote method is defined as serviceName + "." + name.
// Parameters:
//  - prototype         a prototype for properly convert result
//  - name              a name of the command to call.
//  - correlationId     (optional) transaction id to trace execution through call chain.
//  - params            command parameters.
// Retruns: result or error.
func (c *CommandableGrpcClient) CallCommand(prototype reflect.Type, name string, correlationId string, params *cdata.AnyValueMap) (result interface{}, err error) {
	method := c.Name + "." + name
	timing := c.Instrument(correlationId, method)

	var jsonArgs string
	if params != nil {
		jsonRes, err := json.Marshal(params.Value())
		jsonArgs = string(jsonRes)
		if err != nil {
			return nil, err
		}
	}

	request := grpcproto.InvokeRequest{
		Method:        method,
		CorrelationId: correlationId,
		ArgsEmpty:     params == nil,
		ArgsJson:      jsonArgs,
	}

	var response grpcproto.InvokeReply
	err = c.Call("invoke", correlationId, &request, &response)

	timing.EndTiming()

	// Handle unexpected error
	if err != nil {
		return c.InstrumentError(correlationId, method, err, response)
	}

	// Handle error response
	if response.Error != nil {
		var errDesc cerr.ErrorDescription
		errDescJson, _ := json.Marshal(response.Error)
		json.Unmarshal(errDescJson, errDesc)
		err = cerr.ApplicationErrorFactory.Create(&errDesc)
		return nil, err
	}

	// Handle empty response
	if response.ResultEmpty || response.ResultJson == "" {
		return nil, nil
	}

	// Handle regular response
	if prototype != nil {
		return rpcclients.ConvertComandResult([]byte(response.ResultJson), prototype)
	}

	return []byte(response.ResultJson), nil
}
