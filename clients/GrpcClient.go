package clients

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

func NewGrpcClient(name string) *GrpcClient {
	return &GrpcClient{
		name: name,
	}
}

func (c *GrpcClient) Configure(config *config.ConfigParams) {
	host := config.GetAsStringWithDefault("connection.host", "localhost")
	port := config.GetAsStringWithDefault("connection.port", "8090")
	c.address = host + ":" + port
}

func (c *GrpcClient) IsOpen() bool {
	return c.connection != nil
}

func (c *GrpcClient) Open(correlationId string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.connection = conn

	return nil
}

func (c *GrpcClient) Close(correlationId string) error {
	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}

	return nil
}

func (c *GrpcClient) Call(method string, correlationId string, request interface{}, response interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	method = "/" + c.name + "/" + method
	err := c.connection.Invoke(ctx, method, request, response)

	return err
}
