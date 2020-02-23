package services

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
Abstract service that receives remove calls via GRPC protocol.

 Configuration parameters

- dependencies:
  - endpoint:              override for GRPC Endpoint dependency
  - controller:            override for Controller dependency
- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- credential - the HTTPS credentials:
  - ssl_key_file:         the SSL private key in PEM
  - ssl_crt_file:         the SSL certificate in PEM
  - ssl_ca_file:          the certificate authorities (root cerfiticates) in PEM

 References

- \*:logger:\*:\*:1.0               (optional) https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger components to pass log messages
- \*:counters:\*:\*:1.0             (optional) https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters components to pass collected measurements
- \*:discovery:\*:\*:1.0            (optional) https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery services to resolve connection
- \*:endpoint:grpc:\*:1.0           (optional) GrpcEndpoint reference

See GrpcClient

 Example

    class MyGrpcService extends GrpcService {
       private _controller: IMyController;
       ...
       func (c*GrpcService) constructor() {
          base("... path to proto ...", ".. service name ...");
          c.DependencyResolver.put(
              "controller",
              new Descriptor("mygroup","controller","*","*","1.0")
          );
       }

       func (c*GrpcService) setReferences(references: IReferences) {
          base.setReferences(references);
          c._controller = c.DependencyResolver.getRequired<IMyController>("controller");
       }

       func (c*GrpcService) register() {
           registerMethod("get_mydata", null, (call, callback) => {
               let correlationId = call.request.correlationId;
               let id = call.request.id;
               c._controller.getMyData(correlationId, id, callback);
           });
           ...
       }
    }

    let service = new MyGrpcService();
    service.configure(ConfigParams.fromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));
    service.setReferences(References.fromTuples(
       new Descriptor("mygroup","controller","default","default","1.0"), controller
    ));

    service.open("123", (err) => {
       console.log("The GRPC service is running on port 8080");
    });
*/
// implements IOpenable, IConfigurable, IReferenceable,  IUnreferenceable, IRegisterable

type GrpcService struct {
	IRegisterable
	defaultConfig *cconf.ConfigParams
	serviceName   string
	//serviceDescriptor *grpc.ServiceDesc
	config        *cconf.ConfigParams
	references    cref.IReferences
	localEndpoint bool
	//registerable      IRegisterable
	//implementation    interface{}
	opened bool
	//  The GRPC endpoint that exposes c service.
	Endpoint *GrpcEndpoint
	//  The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	//  The logger.
	Logger *clog.CompositeLogger
	//  The performance counters.
	Counters *ccount.CompositeCounters
}

// NewGrpcService creates new instance NewGrpcService
// Parameters:
//  - serviceName string
//  service name from XYZ.pb.go, set "" for use default gRPC commandable protobuf
// Return *GrpcService
func NewGrpcService(serviceName string) *GrpcService {
	gs := GrpcService{}
	gs.serviceName = serviceName
	gs.defaultConfig = cconf.NewConfigParamsFromTuples(
		"dependencies.endpoint", "*:endpoint:grpc:*:1.0",
	)
	gs.DependencyResolver = cref.NewDependencyResolverWithParams(gs.defaultConfig, cref.NewEmptyReferences())
	gs.Logger = clog.NewCompositeLogger()
	gs.Counters = ccount.NewCompositeCounters()
	return &gs
}

/*
  Configures component by passing configuration parameters.

  - config    configuration parameters to be set.
*/
func (c *GrpcService) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config
	c.DependencyResolver.Configure(config)
}

/*
Sets references to dependent components.

- references 	references to locate the component dependencies.
*/
func (c *GrpcService) SetReferences(references cref.IReferences) {
	c.references = references

	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.DependencyResolver.SetReferences(references)

	// Get endpoint
	res := c.DependencyResolver.GetOneOptional("endpoint")
	c.Endpoint, _ = res.(*GrpcEndpoint)
	// Or create a local one
	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.localEndpoint = true
	} else {
		c.localEndpoint = false
	}
	// Add registration callback to the endpoint
	c.Endpoint.Register(c)
}

/*
Unsets (clears) previously set references to dependent components.
*/
func (c *GrpcService) UnsetReferences() {
	// Remove registration callback from endpoint
	if c.Endpoint != nil {
		c.Endpoint.Unregister(c)
		c.Endpoint = nil
	}
}

func (c *GrpcService) createEndpoint() *GrpcEndpoint {
	endpoint := NewGrpcEndpoint()

	if c.config != nil {
		endpoint.Configure(c.config)
	}

	if c.references != nil {
		endpoint.SetReferences(c.references)
	}
	return endpoint
}

/*
Adds instrumentation to log calls and measure call time.
It returns a Timing object that is used to end the time measurement.

- correlationId     (optional) transaction id to trace execution through call chain.
- name              a method name.
Return Timing object to end the time measurement.
*/
func (c *GrpcService) Instrument(correlationId string, name string) *ccount.Timing {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".exec_count")
	return c.Counters.BeginTiming(name + ".exec_time")
}

/*
Adds instrumentation to error handling.

- correlationId     (optional) transaction id to trace execution through call chain.
- name              a method name.
- err               an occured error
- result            (optional) an execution result
- callback          (optional) an execution callback
*/
func (c *GrpcService) InstrumentError(correlationId string, name string, errIn error,
	resIn interface{}) (result interface{}, err error) {
	if errIn != nil {
		c.Logger.Error(correlationId, errIn, "Failed to execute %s method", name)
		c.Counters.IncrementOne(name + ".exec_errors")
	}

	return resIn, errIn
}

/*
Checks if the component is opened.
Return true if the component has been opened and false otherwise.
*/
func (c *GrpcService) IsOpen() bool {
	return c.opened
}

/*
Opens the component.
Parameters:
- correlationId 	(optional) transaction id to trace execution through call chain.
- callback 			callback function that receives error or null no errors occured.
*/
func (c *GrpcService) Open(correlationId string) (err error) {
	if c.opened {
		return nil
	}

	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.Endpoint.Register(c)
		c.localEndpoint = true
	}

	if c.localEndpoint {
		opnErr := c.Endpoint.Open(correlationId)
		if opnErr != nil {
			c.opened = false
			return opnErr
		}
	}
	c.opened = true
	return nil
}

/*
Closes component and frees used resources.
Parameters:
- correlationId 	(optional) transaction id to trace execution through call chain.
- callback 			callback function that receives error or null no errors occured.
*/
func (c *GrpcService) Close(correlationId string) (err error) {
	if !c.opened {
		return nil
	}

	if c.Endpoint == nil {
		return cerr.NewInvalidStateError(correlationId, "NO_Endpoint", "HTTP endpoint is missing")
	}

	if c.localEndpoint {
		clsErr := c.Endpoint.Close(correlationId)
		if clsErr != nil {
			c.opened = false
			return clsErr
		}
	}
	c.opened = false
	return nil

}

/*
Registers a commandable method in c objects GRPC server (service) by the given name.,
Parameters:
- method        the GRPC method name.
- schema        the schema to use for parameter validation.
- action        the action to perform at the given route.
*/
func (c *GrpcService) RegisterCommadableMethod(method string, schema *cvalid.Schema,
	action func(correlationId string, data *crun.Parameters) (result interface{}, err error)) {
	c.Endpoint.RegisterCommadableMethod(method, schema, action)
}
