package services

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	rpccon "github.com/pip-services3-go/pip-services3-rpc-go/connect"

	grpcproto "github.com/pip-services3-go/pip-services3-grpc-go/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
Used for creating GRPC endpoints. An endpoint is a URL, at which a given service can be accessed by a client.

 Configuration parameters

Parameters to pass to the configure method for component configuration:

- connection(s) - the connection resolver"s connections:
    - "connection.discovery_key" - the key to use for connection resolving in a discovery service;
    - "connection.protocol" - the connection"s protocol;
    - "connection.host" - the target host;
    - "connection.port" - the target port;
    - "connection.uri" - the target URI.
- credential - the HTTPS credentials:
    - "credential.ssl_key_file" - the SSL private key in PEM
    - "credential.ssl_crt_file" - the SSL certificate in PEM
    - "credential.ssl_ca_file" - the certificate authorities (root cerfiticates) in PEM

 References

A logger, counters, and a connection resolver can be referenced by passing the
following references to the object"s setReferences method:

- logger: "\*:logger:\*:\*:1.0";
- counters: "\*:counters:\*:\*:1.0";
- discovery: "\*:discovery:\*:\*:1.0" (for the connection resolver).

 Examples

    func (c*GrpcEndpoint) MyMethod(_config: ConfigParams, _references: IReferences) {
        let endpoint = new HttpEndpoint();
        if (c._config)
            endpoint.configure(c._config);
        if (c._references)
            endpoint.setReferences(c._references);
        ...

        c._endpoint.open(correlationId, (err) => {
                c._opened = err == nil;
                callback(err);
            });
        ...
    }
*/
// IOpenable, IConfigurable, IReferenceable

type GrpcEndpoint struct {
	defaultConfig *cconf.ConfigParams

	server *grpc.Server

	connectionResolver *rpccon.HttpConnectionResolver
	logger             *clog.CompositeLogger
	counters           *ccount.CompositeCounters
	maintenanceEnabled bool
	fileMaxSize        int64
	uri                string
	registrations      []IRegisterable
	commandableMethods map[string]func(correlationId string, args *crun.Parameters) (result interface{}, err error)
	commandableSchemas map[string]*cvalid.Schema
	interceptors       []grpc.ServerOption
}

func NewGrpcEndpoint() *GrpcEndpoint {
	ge := GrpcEndpoint{}
	ge.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", 3000,

		"credential.ssl_key_file", "",
		"credential.ssl_crt_file", "",
		"credential.ssl_ca_file", "",

		"options.maintenance_enabled", false,
		"options.request_max_size", 1024*1024,
		"options.file_max_size", 200*1024*1024,
		"options.connect_timeout", 60000,
		"options.debug", true,
	)

	ge.connectionResolver = rpccon.NewHttpConnectionResolver()
	ge.logger = clog.NewCompositeLogger()
	ge.counters = ccount.NewCompositeCounters()
	ge.maintenanceEnabled = false
	ge.fileMaxSize = 200 * 1024 * 1024
	ge.registrations = make([]IRegisterable, 0)
	ge.commandableMethods = make(map[string]func(correlationId string, args *crun.Parameters) (result interface{}, err error), 0)
	ge.commandableSchemas = make(map[string]*cvalid.Schema, 0)
	ge.interceptors = make([]grpc.ServerOption, 0, 0)
	return &ge
}

/*
   Configures c HttpEndpoint using the given configuration parameters.

   Configuration parameters:
   - connection(s) - the connection resolver"s connections;
       - "connection.discovery_key" - the key to use for connection resolving in a discovery service;
       - "connection.protocol" - the connection"s protocol;
       - "connection.host" - the target host;
       - "connection.port" - the target port;
       - "connection.uri" - the target URI.
       - "credential.ssl_key_file" - SSL private key in PEM
       - "credential.ssl_crt_file" - SSL certificate in PEM
       - "credential.ssl_ca_file" - Certificate authority (root certificate) in PEM

   - config    configuration parameters, containing a "connection(s)" section.

   See https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/classes/config.configparams.html ConfigParams (in the PipServices "Commons" package)
*/
func (c *GrpcEndpoint) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.connectionResolver.Configure(config)

	c.maintenanceEnabled = config.GetAsBooleanWithDefault("options.maintenance_enabled", c.maintenanceEnabled)
	c.fileMaxSize = config.GetAsLongWithDefault("options.file_max_size", c.fileMaxSize)
}

/*
   Sets references to c endpoint"s logger, counters, and connection resolver.

   References:
   - logger: "*:logger:*:*:1.0"
   - counters: "*:counters:*:*:1.0"
   - discovery: "*:discovery:*:*:1.0" (for the connection resolver)

   - references    an IReferences object, containing references to a logger, counters,
                        and a connection resolver.

   See IReferences (in the PipServices "Commons" package)
*/
func (c *GrpcEndpoint) SetReferences(references cref.IReferences) {
	c.logger.SetReferences(references)
	c.counters.SetReferences(references)
	c.connectionResolver.SetReferences(references)
}

/*
   Return whether or not c endpoint is open with an actively listening GRPC server.
*/
func (c *GrpcEndpoint) IsOpen() bool {
	return c.server != nil
}

/*
   Opens a connection using the parameters resolved by the referenced connection
   resolver and creates a GRPC server (service) using the set options and parameters.

   - correlationId     (optional) transaction id to trace execution through call chain.
   - callback          (optional) the function to call once the opening process is complete.
                            Will be called with an error if one is raised.
*/
func (c *GrpcEndpoint) Open(correlationId string) (err error) {
	if c.IsOpen() {
		return nil
	}

	connection, credential, err := c.connectionResolver.Resolve(correlationId)
	if err != nil {
		return err
	}

	c.uri = connection.Host() + ":" + strconv.FormatInt(int64(connection.Port()), 16)

	opts := []grpc.ServerOption{}

	if len(c.interceptors) > 0 {
		// Add interceptors
		opts = append(opts, c.interceptors...)
	}

	if connection.Protocol() == "https" {
		sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")
		creds, _ := credentials.NewServerTLSFromFile(sslCrtFile, sslKeyFile)
		opts = append(opts, grpc.Creds(creds))
	}

	lis, _ := net.Listen("tcp", c.uri)

	// Create instance of express application

	c.server = grpc.NewServer(opts...)

	err = c.connectionResolver.Register(correlationId)
	if err != nil {
		return nil
	}
	c.logger.Debug(correlationId, "Opened GRPC service at tcp:\\\\%s", c.uri)

	// Start operations
	c.performRegistrations()

	go func() {
		servErr := c.server.Serve(lis)
		if servErr != nil {
			err := cerr.NewConnectionError(correlationId, "CANNOT_CONNECT", "Opening GRPC service failed").
				Wrap(err).WithDetails("url", c.uri)
			c.logger.Debug(correlationId, "Opened GRPC service at %s", err)
		}
	}()

	return nil
}

/*
Closes c endpoint and the GRPC server (service) that was opened earlier.

- correlationId     (optional) transaction id to trace execution through call chain.
- callback          (optional) the function to call once the closing process is complete.
                         Will be called with an error if one is raised.
*/
func (c *GrpcEndpoint) Close(correlationId string) (err error) {
	if c.server != nil {
		c.uri = ""

		c.commandableMethods = nil
		c.commandableSchemas = nil

		c.server.GracefulStop()
		c.logger.Debug(correlationId, "Closed GRPC service at %s", c.uri)
		c.server = nil
	}
	return nil

}

// GetServer return working gRPC server for register services
// Retruns *grpc.Server
func (c *GrpcEndpoint) GetServer() *grpc.Server {
	return c.server
}

// AddUnaryInterceptor Registers a middleware for methods in GRPC endpoint.
// See https://github.com/grpc/grpc-go/tree/master/examples/features/interceptor
// Parameters:
// - interceptors ...grpc.ServerOption
// interceptor functions (Stream or Unary use grpc.UnaryInterceptor() or grpc.StreamInterceptor() for inflate in grpc.ServerOption)
func (c *GrpcEndpoint) AddInterceptors(interceptors ...grpc.ServerOption) {
	c.interceptors = append(c.interceptors, interceptors...)
}

/*
   Registers a registerable object for dynamic endpoint discovery.

   - registration      the registration to add.

   See IRegisterable
*/
func (c *GrpcEndpoint) Register(registration IRegisterable) {
	c.registrations = append(c.registrations, registration)
}

/*
   Unregisters a registerable object, so that it is no longer used in dynamic
   endpoint discovery.

   - registration      the registration to remove.

   See IRegisterable
*/
func (c *GrpcEndpoint) Unregister(registration IRegisterable) {
	for i := 0; i < len(c.registrations); {
		if c.registrations[i] == registration {
			if i == len(c.registrations)-1 {
				c.registrations = c.registrations[:i]
			} else {
				c.registrations = append(c.registrations[:i], c.registrations[i+1:]...)
			}
		} else {
			i++
		}
	}
}

func (c *GrpcEndpoint) performRegistrations() {
	for _, registration := range c.registrations {
		registration.Register()
	}
	c.registerCommandableService()
}

func (c *GrpcEndpoint) registerCommandableService() {
	if len(c.commandableMethods) == 0 {
		return
	}
	invokeMediator := InvokeComandMediator{InvokeFunc: c.invoke}
	grpcproto.RegisterCommandableServer(c.server, &invokeMediator)
}

/*
   Registers a service with related implementation
   - implementation the service implementation method Invoke.
*/
func (c *GrpcEndpoint) RegisterService(sd *grpc.ServiceDesc, implementation interface{}) {
	if c.server != nil {
		c.server.RegisterService(sd, implementation)
	}
}

// Invoke method for implements interface grpcproto.CommandableServer
func (c *GrpcEndpoint) invoke(ctx context.Context, request *grpcproto.InvokeRequest) (response *grpcproto.InvokeReply, err error) {

	method := request.Method

	var action func(correlationId string, args *crun.Parameters) (result interface{}, err error)
	if len(c.commandableMethods) > 0 {
		action = c.commandableMethods[method]
	}

	correlationId := request.CorrelationId

	// Handle method not found
	if action == nil {
		appErr := cerr.NewInvocationError(correlationId, "METHOD_NOT_FOUND", "Method "+method+" was not found").
			WithDetails("method", method)

		var errDesc grpcproto.ErrorDescription
		errDescJson, _ := json.Marshal(appErr)
		json.Unmarshal(errDescJson, errDesc)

		response = &grpcproto.InvokeReply{
			Error:       &errDesc,
			ResultEmpty: true,
			ResultJson:  "",
		}
		return response, nil
	}
	// Convert arguments
	argsEmpty := request.ArgsEmpty
	argsJson := request.ArgsJson
	var args *crun.Parameters = crun.NewEmptyParameters()

	if !argsEmpty && argsJson != "" {
		var buf map[string]interface{}
		err := json.Unmarshal([]byte(argsJson), buf)
		if err == nil {
			args.Append(buf)
		}
	}

	// Todo: Validate schema
	// schema := c.commandableSchemas[method];
	// if schema != nil {
	//     //...
	// }

	// Call command action
	result, err := action(correlationId, args)
	// Process result and generate response
	if err != nil {

		appErr := cerr.ErrorDescriptionFactory.Create(err)
		var errDesc grpcproto.ErrorDescription
		errDescJson, _ := json.Marshal(appErr)
		json.Unmarshal(errDescJson, errDesc)

		response = &grpcproto.InvokeReply{
			Error:       &errDesc,
			ResultEmpty: true,
			ResultJson:  "",
		}
	} else {
		resJson, _ := json.Marshal(result)
		response = &grpcproto.InvokeReply{
			Error:       nil,
			ResultEmpty: result == nil,
			ResultJson:  string(resJson),
		}
	}

	return response, err
}

/*
Registers a commandable method in c objects GRPC server (service) by the given name.,

- method        the GRPC method name.
- schema        the schema to use for parameter validation.
- action        the action to perform at the given route.
*/
func (c *GrpcEndpoint) RegisterCommadableMethod(method string, schema *cvalid.Schema,
	action func(correlationId string, args *crun.Parameters) (result interface{}, err error)) {

	if c.commandableMethods == nil {
		c.commandableMethods = make(map[string]func(correlationId string, args *crun.Parameters) (result interface{}, err error))
	}
	c.commandableMethods[method] = action

	if c.commandableSchemas == nil {
		c.commandableSchemas = make(map[string]*cvalid.Schema)
	}
	c.commandableSchemas[method] = schema
}
