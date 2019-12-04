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
    func  (chc * MyCommandableHttpClient) GetData(correlationId string, id string) (res interface{}, err error) {

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

	"github.com/pip-services3-go/pip-services3-commons-go/config"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	name       string
	address    string
	connection *grpc.ClientConn
}

// Creates a new instance of the client.
// Parameters:
// 			- baseRoute string
// 			a base route for remote service.
// Returns *GrpcClient
func NewGrpcClient(name string) *GrpcClient {
	return &GrpcClient{
		name: name,
	}
}

// Configures component by passing configuration parameters.
// Parameters:
// 			- config *config.ConfigParams
// 			onfiguration parameters to be set.
func (c *GrpcClient) Configure(config *config.ConfigParams) {
	host := config.GetAsStringWithDefault("connection.host", "localhost")
	port := config.GetAsStringWithDefault("connection.port", "8090")
	c.address = host + ":" + port
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
	// Set up a connection to the server.
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.connection = conn

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

// Calls a remote method via HTTP/REST protocol.
// Parameters:
// 		- method string
// 		HTTP method: "get", "head", "post", "put", "delete"
// 		- route string
// 		a command route. Base route will be added to this route
// 		- correlationId string
// 		transaction id to trace execution through call chain.
// 		- request interface{}
// 		query parameters.
// 		- response interface{}
// 		 responce body object.
// Returns error
func (c *GrpcClient) Call(method string, correlationId string, request interface{}, response interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	method = "/" + c.name + "/" + method
	err := c.connection.Invoke(ctx, method, request, response)

	return err
}
