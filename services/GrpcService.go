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
	defaultConfig  *cconf.ConfigParams
	serviceName    string
	config         *cconf.ConfigParams
	references     cref.IReferences
	localEndpoint  bool
	registerable   IRegisterable
	implementation interface{}
	interceptors   []interface{}
	opened         bool

	/*
	  The GRPC endpoint that exposes c service.
	*/
	Endpoint *GrpcEndpoint
	/*
	  The dependency resolver.
	*/
	DependencyResolver *cref.DependencyResolver
	/*
	  The logger.
	*/
	Logger *clog.CompositeLogger
	/*
	  The performance counters.
	*/
	Counters *ccount.CompositeCounters
}

func NewGrpcService(serviceName string) *GrpcService {
	gs := GrpcService{}
	gs.serviceName = serviceName
	gs.defaultConfig = cconf.NewConfigParamsFromTuples(
		"dependencies.endpoint", "*:endpoint:grpc:*:1.0",
	)
	gs.DependencyResolver = cref.NewDependencyResolverWithParams(gs.defaultConfig, cref.NewEmptyReferences())
	gs.Logger = clog.NewCompositeLogger()
	gs.Counters = ccount.NewCompositeCounters()

	gs.interceptors = make([]interface{}, 0)

	// gs.registerable = {
	//     register: () => {
	//         c.registerService();
	//     }
	// }
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
func (c *GrpcService) setReferences(references cref.IReferences) {
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
	c.Endpoint.Register(c.registerable)
}

/*
Unsets (clears) previously set references to dependent components.
*/
func (c *GrpcService) unsetReferences() {
	// Remove registration callback from endpoint
	if c.Endpoint != nil {
		c.Endpoint.Unregister(c.registerable)
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
@returns Timing object to end the time measurement.
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

@returns true if the component has been opened and false otherwise.
*/
func (c *GrpcService) IsOpen() bool {
	return c.opened
}

/*
Opens the component.

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

// private registerService() {
//     // Register implementations
//     c.implementation = {};
//     c.interceptors = [];
//     c.register();

//     // Load service
//     let grpc = require("grpc");
//     let service = c.service;

//     // Dynamically load service
//     if (service == null && _.isString(c.protoPath)) {
//         let protoLoader = require("@grpc/proto-loader");

//         let options = c.packageOptions || {
//             keepCase: true,
//             longs: Number,
//             enums: Number,
//             defaults: true,
//             oneofs: true
//         };

//         let packageDefinition = protoLoader.loadSync(c.protoPath, options);
//         let packageObject = grpc.loadPackageDefinition(packageDefinition);
//         service = c.getServiceByName(packageObject, c.serviceName);
//     }
//     // Statically load service
//     else {
//         service = c.getServiceByName(c.service, c.serviceName);
//     }

//     // Register service if it is set
//     if (service) {
//         c.Endpoint.registerService(service, c.implementation);
//     }
// }

/*
Registers a method in GRPC service.

- name          a method name
- schema        a validation schema to validate received parameters.
- action        an action function that is called when operation is invoked.
*/
//  func (c *GrpcService)  RegisterMethod(name string, schema *cvalid.Schema,
//     action func(call: any, callback: (err: any, message: any) => void) => void) {
//     if (c.implementation == null) return;

//     c.implementation[name] = (call, callback) => {
//         async.each(c.interceptors, (interceptor, cb) => {
//             interceptor(call, callback, cb);
//         }, (err) => {
//             // Validate object
//             if (schema && call && call.request) {
//                 let value = call.request;
//                 if (_.isFunction(value.toObject))
//                     value = value.toObject();

//                 // Perform validation
//                 let correlationId = value.correlation_id;
//                 let err = schema.validateAndReturnException(correlationId, value, false);
//                 if (err) {
//                     callback(err, null);
//                     return;
//                 }
//             }

//             action.call(c, call, callback);
//         });
//     };
// }

/*
Registers a commandable method in c objects GRPC server (service) by the given name.,

- method        the GRPC method name.
- schema        the schema to use for parameter validation.
- action        the action to perform at the given route.
*/
func (c *GrpcService) RegisterCommadableMethod(method string, schema *cvalid.Schema,
	action func(correlationId string, data *crun.Parameters) (result interface{}, err error)) {

	c.Endpoint.RegisterCommadableMethod(method, schema, action)
}

/*
Registers a middleware for methods in GRPC endpoint.

- action        an action function that is called when middleware is invoked.
*/
//  func (c *GrpcService) RegisterInterceptor(
//     action func(call: any, callback: (err: any, result: any) => void, next: () => void) => void) {
//     if (c.Endpoint == null) return;

//     c.interceptors.push((call, callback, next) => {
//         action.call(c, call, callback, next);
//     });
// }

// /*
// Registers all service routes in gRPC endpoint.
// c method is called by the service and must be overriden
// in child classes.
//  */
// func (c*GrpcService) abstract register();
