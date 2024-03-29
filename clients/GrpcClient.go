package clients

import (
	"context"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"
	rpccon "github.com/pip-services3-go/pip-services3-rpc-go/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

/*
GrpcClient abstract client that calls commandable HTTP service.

Commandable services are generated automatically for ICommandable objects. Each command is exposed as POST operation that receives all parameters in body object.

Configuration parameters:

  base_route: base route for remote URI
  connection(s):
  discovery_key: (optional) a key to retrieve the connection from IDiscovery
  protocol: connection protocol: http or https
  host: host name or IP address
  port: port number
  uri: resource URI or connection string with all parameters in it
  options:
  retries: number of retries (default: 3)
  connect_timeout: connection timeout in milliseconds (default: 10 sec)
  timeout: invocation timeout in milliseconds (default: 10 sec)

References:

*:logger:*:*:1.0 (optional) ILogger components to pass log messages
*:counters:*:*:1.0 (optional) ICounters components to pass collected measurements
*:discovery:*:*:1.0 (optional) IDiscovery services to resolve connection

Example:

type MyCommandableHttpClient struct{
 	*CommandableHttpClient
}
    func  (c *MyCommandableHttpClient) GetData(correlationId string, id string) (res interface{}, err error) {

        req := &testproto.MyDataIdRequest{
            CorrelationId: correlationId,
            mydataId:       id,
        }

        reply := new(testproto.MyData)
        err = c.Call("get_mydata_by_id", correlationId, req, reply)
        c.Instrument(correlationId, "mydata.get_one_by_id")
        if err != nil {
            return nil, err
        }
        result = toMyData(reply)
        if result != nil && result.Id == "" && result.Key == "" {
            result = nil
        }
        return result, nil
	}

var client = NewMyCommandableHttpClient();
client.Configure(NewConfigParamsFromTuples(
    "connection.protocol", "http",
    "connection.host", "localhost",
    "connection.port", 8080,
));

result, err := client.GetData("123", "1")
...
*/
type GrpcClient struct {
	address string
	name    string

	defaultConfig *cconf.ConfigParams
	//	The GRPC client.
	Client grpcproto.CommandableClient
	// The GRPC connection
	Connection *grpc.ClientConn
	//	The connection resolver.
	ConnectionResolver *rpccon.HttpConnectionResolver
	//	The logger.
	Logger *clog.CompositeLogger
	//	The performance counters.
	Counters *ccount.CompositeCounters
	//	The configuration options.
	Options *cconf.ConfigParams
	//	The connection timeout in milliseconds.
	ConnectTimeout time.Duration
	//	The invocation timeout in milliseconds.
	Timeout time.Duration
	//	The remote service uri which is calculated on open.
	Uri string
	// interceptors
	interceptors []grpc.DialOption
}

// NewGrpcClient method are creates a new instance of the client.
// Parameters:
//   - baseRoute string
//   a base route for remote service.
// Returns *GrpcClient
func NewGrpcClient(name string) *GrpcClient {
	c := GrpcClient{
		name: name,
	}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "8090",

		"options.connect_timeout", 10000,
		"options.timeout", 10000,
		"options.retries", 3,
		"options.debug", true,
	)
	c.ConnectionResolver = rpccon.NewHttpConnectionResolver()
	c.Logger = clog.NewCompositeLogger()
	c.Counters = ccount.NewCompositeCounters()
	c.Options = cconf.NewEmptyConfigParams()
	c.ConnectTimeout = 10000 * time.Millisecond
	c.Timeout = 10000 * time.Millisecond
	c.interceptors = make([]grpc.DialOption, 0, 0)
	return &c
}

// Configure method are configures component by passing configuration parameters.
// Parameters:
//   - config *config.ConfigParams
//   configuration parameters to be set.
func (c *GrpcClient) Configure(config *cconf.ConfigParams) {
	host := config.GetAsStringWithDefault("connection.host", "localhost")
	port := config.GetAsStringWithDefault("connection.port", "8090")

	c.ConnectTimeout = time.Duration(config.GetAsIntegerWithDefault("connection.connect_timeout", 10000)) * time.Millisecond
	c.Timeout = time.Duration(config.GetAsIntegerWithDefault("connection.timeout", 10000)) * time.Millisecond
	c.ConnectionResolver.Configure(config)
	c.address = host + ":" + port
}

// SetReferences method are sets references to dependent components.
//   - references  cref.IReferences
//   references to locate the component dependencies.
func (c *GrpcClient) SetReferences(references cref.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.ConnectionResolver.SetReferences(references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - name              a method name.
// Returns: Timing object to end the time measurement.
func (c *GrpcClient) Instrument(correlationId string, name string) *ccount.CounterTiming {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".call_count")
	return c.Counters.BeginTiming(name + ".call_time")
}

// InstrumentError mrthod are adds instrumentation to error handling.
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - name              a method name.
//   - err               an occured error
//   - result            (optional) an execution result
// Retruns: result interface{}, err error
// input result and error.
func (c *GrpcClient) InstrumentError(correlationId string, name string, inErr error, inRes interface{}) (result interface{}, err error) {
	if inErr != nil {
		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
		c.Counters.IncrementOne(name + ".call_errors")
	}

	return inRes, inErr
}

// IsOpen method are checks if the component is opened.
// Returns bool
// true if the component has been opened and false otherwise.
func (c *GrpcClient) IsOpen() bool {
	return c.Connection != nil
}

// AddInterceptors method are registers a middleware for methods in gRPC client.
// See https://github.com/grpc/grpc-go/tree/master/examples/features/interceptor
// Parameters:
//   - interceptors ...grpc.DialOption
// interceptor functions (Stream or Unary use grpc.WithUnaryInterceptor() or grpc.WithStreamInterceptor() for inflate in grpc.ServerOption)
func (c *GrpcClient) AddInterceptors(interceptors ...grpc.DialOption) {
	c.interceptors = append(c.interceptors, interceptors...)
}

// Open method are opens the component.
// Parameters:
//   - correlationId string
//   transaction id to trace execution through call chain.
// Returns error
// error or nil
func (c *GrpcClient) Open(correlationId string) error {

	if c.IsOpen() {
		return nil
	}
	connection, credential, err := c.ConnectionResolver.Resolve(correlationId)
	if err != nil {
		return err
	}
	c.Uri = connection.Uri()
	// Set up a connection to the server.
	opts := []grpc.DialOption{
		grpc.WithTimeout(c.ConnectTimeout),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Timeout: c.Timeout}),
	}

	if len(c.interceptors) > 0 {
		// Add interceptors
		opts = append(opts, c.interceptors...)
	}

	if connection.Protocol() == "https" {
		//sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")
		transport, err := credentials.NewClientTLSFromFile(sslCrtFile, c.name)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.WithTransportCredentials(transport))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(c.address, opts...)
	if err != nil {
		return err
	}
	c.Connection = conn
	c.Client = grpcproto.NewCommandableClient(conn)
	return nil
}

// Close method are closes component and frees used resources.
// Parameters:
//   - correlationId string
//   transaction id to trace execution through call chain.
// Returns error
func (c *GrpcClient) Close(correlationId string) error {
	if c.Connection != nil {
		c.Connection.Close()
		c.Connection = nil
	}
	return nil
}

// Call method are calls a remote method via gRPC protocol.
// Parameters:
//   - method string
//   gRPC method name
//   - correlationId string
//   transaction id to trace execution through call chain.
//   - request interface{}
//    request query parameters.
//   - response interface{}
//   - response body object.
// Returns error
func (c *GrpcClient) Call(method string, correlationId string, request interface{}, response interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	method = "/" + c.name + "/" + method
	err := c.Connection.Invoke(ctx, method, request, response)
	return err
}

// CallWithContext method are calls a remote method via gRPC protocol.
// Parameters:
//   - ctx context
//   - correlationId string
//   transaction id to trace execution through call chain.
//   - method string//   gRPC method name
//   - request interface{}
//    request query parameters.
//   - response interface{}
//   - response body object.
// Returns error
func (c *GrpcClient) CallWithContext(ctx context.Context, correlationId string, method string, request interface{}, response interface{}) error {
	method = "/" + c.name + "/" + method
	err := c.Connection.Invoke(ctx, method, request, response)
	return err
}

// AddFilterParams method are adds filter parameters (with the same name as they defined)
// to invocation parameter map.
//   - params        invocation parameters.
//   - filter        (optional) filter parameters
// Return invocation parameters with added filter parameters.
func (c *GrpcClient) AddFilterParams(params *cdata.StringValueMap, filter *cdata.FilterParams) *cdata.StringValueMap {

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	if filter != nil {
		for k, v := range filter.Value() {
			params.Put(k, v)
		}
	}
	return params
}

// AddPagingParams method are adds paging parameters (skip, take, total) to invocation parameter map.
//   - params        invocation parameters.
//   - paging        (optional) paging parameters
// Return invocation parameters with added paging parameters.
func (c *GrpcClient) AddPagingParams(params *cdata.StringValueMap, paging *cdata.PagingParams) *cdata.StringValueMap {
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if paging != nil {
		params.Put("total", paging.Total)
		if paging.Skip != nil {
			params.Put("skip", *paging.Skip)
		}
		if paging.Take != nil {
			params.Put("take", *paging.Take)
		}
	}
	return params
}
