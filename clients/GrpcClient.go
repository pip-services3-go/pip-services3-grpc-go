package clients

/*
Abstract client that calls commandable HTTP service.

Commandable services are generated automatically for ICommandable objects. Each command is exposed as POST operation that receives all parameters in body object.

Configuration parameters
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
References
*:logger:*:*:1.0 (optional) ILogger components to pass log messages
*:counters:*:*:1.0 (optional) ICounters components to pass collected measurements
*:discovery:*:*:1.0 (optional) IDiscovery services to resolve connection

Example
type MyCommandableHttpClient {
CommandableHttpClient

}
    func  (c *MyCommandableHttpClient) GetData(correlationId string, id string) (res interface{}, err error) {

       res, err = chc.callCommand(
           "get_data",
           correlationId,
           { id: id });
	}

var client = NewMyCommandableHttpClient();
client.Configure(NewConfigParamsFromTuples(
    "connection.protocol", "http",
    "connection.host", "localhost",
    "connection.port", 8080
));

client.GetData("123", "1", (err, result) => {
...
});
*/
import (
	"context"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"
	rpccon "github.com/pip-services3-go/pip-services3-rpc-go/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	address string
	name    string

	defaultConfig *cconf.ConfigParams
	//	The GRPC client.
	Client grpcproto.CommandableClient
	// The GRPC connection
	connection *grpc.ClientConn
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
}

// Creates a new instance of the client.
// Parameters:
// 			- baseRoute string
// 			a base route for remote service.
// Returns *GrpcClient
func NewGrpcClient(name string) *GrpcClient {
	gc := GrpcClient{
		name: name,
	}
	gc.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", 3000,

		"options.connect_timeout", 10000,
		"options.timeout", 10000,
		"options.retries", 3,
		"options.debug", true,
	)
	gc.ConnectionResolver = rpccon.NewHttpConnectionResolver()
	gc.Logger = clog.NewCompositeLogger()
	gc.Counters = ccount.NewCompositeCounters()
	gc.Options = cconf.NewEmptyConfigParams()
	gc.ConnectTimeout = 10000 * time.Millisecond
	gc.Timeout = 10000 * time.Millisecond

	return &gc
}

// Configures component by passing configuration parameters.
// Parameters:
// 			- config *config.ConfigParams
// 			onfiguration parameters to be set.
func (c *GrpcClient) Configure(config *cconf.ConfigParams) {
	host := config.GetAsStringWithDefault("connection.host", "localhost")
	port := config.GetAsStringWithDefault("connection.port", "8090")

	c.ConnectTimeout = time.Duration(config.GetAsIntegerWithDefault("connection.connect_timeout", 10000)) * time.Millisecond
	c.Timeout = time.Duration(config.GetAsIntegerWithDefault("connection.timeout", 10000)) * time.Millisecond

	c.address = host + ":" + port
}

// Sets references to dependent components.
// - references  cref.IReferences
//	references to locate the component dependencies.
func (c *GrpcClient) SetReferences(references cref.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.ConnectionResolver.SetReferences(references)
}

/*
  Adds instrumentation to log calls and measure call time.
  It returns a Timing object that is used to end the time measurement.

  - correlationId     (optional) transaction id to trace execution through call chain.
  - name              a method name.
  @returns Timing object to end the time measurement.
*/
func (c *GrpcClient) Instrument(correlationId string, name string) *ccount.Timing {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".call_count")
	return c.Counters.BeginTiming(name + ".call_time")
}

/*
  Adds instrumentation to error handling.

  - correlationId     (optional) transaction id to trace execution through call chain.
  - name              a method name.
  - err               an occured error
  - result            (optional) an execution result
  - callback          (optional) an execution callback
*/
func (c *GrpcClient) InstrumentError(correlationId string, name string, inErr error, inRes interface{}) (result interface{}, err error) {
	if inErr != nil {
		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
		c.Counters.IncrementOne(name + ".call_errors")
	}

	return inRes, inErr
}

// Checks if the component is opened.
// Returns bool
// true if the component has been opened and false otherwise.
func (c *GrpcClient) IsOpen() bool {
	return c.connection != nil
}

// Opens the component.
// Parameters:
// 			-correlationId string
// 			 transaction id to trace execution through call chain.
// Returns error
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
	c.connection = conn

	c.Client = grpcproto.NewCommandableClient(conn)

	return nil
}

// Closes component and frees used resources.
// Parameters:
// 			- correlationId string
// 			transaction id to trace execution through call chain.
// Returns error
func (c *GrpcClient) Close(correlationId string) error {
	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}
	return nil
}

// Calls a remote method via gRPC protocol.
// Parameters:
// 		- method string
// 		gRPC method name
// 		- correlationId string
// 		transaction id to trace execution through call chain.
// 		- request interface{}
// 		 request query parameters.
// 		- response interface{}
// 		- response body object.
// Returns error
func (c *GrpcClient) Call(method string, correlationId string, request interface{}, response interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	method = "/" + c.name + "/" + method
	err := c.connection.Invoke(ctx, method, request, response)
	return err

}
