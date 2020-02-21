package clients

import (
	"encoding/json"

	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"
)

// import (
// 	pb "github.com/pip-services3-go/pip-services3-grpc-go/protos"
// )

/*
Abstract client that calls commandable GRPC service.

Commandable services are generated automatically for https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/interfaces/commands.icommandable.html ICommandable objects.
Each command is exposed as Invoke method that receives all parameters as args.

 Configuration parameters

- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- options:
  - retries:               number of retries (default: 3)
  - connect_timeout:       connection timeout in milliseconds (default: 10 sec)
  - timeout:               invocation timeout in milliseconds (default: 10 sec)

 References

- *:logger:*:*:1.0         (optional) https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger components to pass log messages
- *:counters:*:*:1.0         (optional) https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters components to pass collected measurements
- *:discovery:*:*:1.0        (optional) https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery services to resolve connection

 Example

    class MyCommandableGrpcClient extends CommandableGrpcClient implements IMyClient {
       ...

        public getData(correlationId: string, id: string,
           callback: (err: any, result: MyData) => void): void {

           c.callCommand(
               "get_data",
               correlationId,
               { id: id },
               (err, result) => {
                   callback(err, result);
               }
            );
        }
        ...
    }

    let client = new MyCommandableGrpcClient();
    client.configure(ConfigParams.fromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));

    client.getData("123", "1", (err, result) => {
    ...
    });
*/
type CommandableGrpcClient struct {
	*GrpcClient
	//The service name
	Name string
}

/*
   Creates a new instance of the client.
   - name     a service name.
*/
func NewCommandableGrpcClient(name string) *CommandableGrpcClient {
	cgc := CommandableGrpcClient{}
	cgc.GrpcClient = NewGrpcClient(name)
	cgc.Name = name
	return &cgc
}

/*
Calls a remote method via GRPC commadable protocol.
The call is made via Invoke method and all parameters are sent in args object.
The complete route to remote method is defined as serviceName + "." + name.

- name              a name of the command to call.
- correlationId     (optional) transaction id to trace execution through call chain.
- params            command parameters.
- callback          callback function that receives result or error.
*/
func (c *CommandableGrpcClient) CallCommand(name string, correlationId string, params interface{}) (result interface{}, err error) {
	method := c.Name + "." + name
	timing := c.Instrument(correlationId, method)

	var jsonArgs string
	if params != nil {
		jsonRes, err := json.Marshal(params)
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
	err = c.Call("invoke", correlationId, request, &response)

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
	var resObj interface{}
	err = json.Unmarshal([]byte(response.ResultJson), &resObj)

	if err != nil {
		return nil, err
	}
	return &resObj, nil

}
