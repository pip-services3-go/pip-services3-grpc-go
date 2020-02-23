package services

import (
	ccomands "github.com/pip-services3-go/pip-services3-commons-go/commands"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
)

/*
Abstract service that receives commands via GRPC protocol
to operations automatically generated for commands defined in https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/interfaces/commands.icommandable.html ICommandable components.
Each command is exposed as invoke method that receives command name and parameters.

Commandable services require only 3 lines of code to implement a robust external
GRPC-based remote interface.

 Configuration parameters

- dependencies:
  - endpoint:              override for HTTP Endpoint dependency
  - controller:            override for Controller dependency
- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from  IDiscovery
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it

 References

- *:logger:\*:\*:1.0               (optional)  ILogger components to pass log messages
- *:counters:\*:\*:1.0             (optional) ICounters components to pass collected measurements
- *:discovery:\*:\*:1.0            (optional)  IDiscovery services to resolve connection
- *:endpoint:grpc:\*:1.0          (optional) GrpcEndpoint reference

See CommandableGrpcClient
See GrpcService

 Example

    class MyCommandableGrpcService extends CommandableGrpcService {
       func (c *CommandableGrpcService ) constructor() {
          base();
          c._dependencyResolver.put(
              "controller",
              new Descriptor("mygroup","controller","*","*","1.0")
          );
       }
    }

    let service = new MyCommandableGrpcService();
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
type CommandableGrpcService struct {
	*GrpcService
	name       string
	commandSet *ccomands.CommandSet
}

/*
   Creates a new instance of the service.
   - name a service name.
*/
func NewCommandableGrpcService(name string) *CommandableGrpcService {
	cgs := CommandableGrpcService{}
	cgs.GrpcService = NewGrpcService("")
	cgs.GrpcService.IRegisterable = &cgs
	cgs.name = name
	cgs.DependencyResolver.Put("controller", "none")
	return &cgs
}

/*
   Registers all service command in gRPC endpoint.
*/
func (c *CommandableGrpcService) Register() {

	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		return
	}
	controller, ok := resCtrl.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error("CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}
	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	var index = 0
	for index = 0; index < len(commands); index++ {
		command := commands[index]

		method := c.name + "." + command.Name()

		c.RegisterCommadableMethod(method, nil,
			func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
				timing := c.Instrument(correlationId, method)
				res, err := command.Execute(correlationId, args)
				timing.EndTiming()
				return c.InstrumentError(correlationId, method, err, res)
			})
	}
}
