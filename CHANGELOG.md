# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> GRPC Component for Pip.Services in Go Changelog

## <a name="1.2.0"></a> 1.2.0 (2021-04-22) 

### Features
* **test** Added TestGrpcClient
* **test** Added TestCommandableGrpcClient

## <a name="1.1.1"></a> 1.1.1 (2021-04-16) 

### Bug Fixes
* Make public connection field in GrpcClient
* Add CallWithContext in GrpcClient 

## <a name="1.1.0"></a> 1.1.0 (2021-04-04) 

### Breaking Changes
* Introduced IGrpcServiceOverrides
* Changed signature NewGrpcService to InheritGrpcService
* Changed signature NewCommandableGrpcService to InheritGrpcService

## <a name="1.0.1-1.0.2"></a> 1.0.1-1.0.2 (2020-11-12) 

### Bug Fixes
* Fix CallCommand method in CommandableGrpcClient


## <a name="1.0.1-1.0.2"></a> 1.0.1-1.0.2 (2020-11-05) 

### Bug Fixes
* Fix default factory
* Fix GRPC Endpoint


## <a name="1.0.0"></a> 1.0.0 (2020-03-05) 

Initial public release

### Features
* **build** factories for creating gRPC services
* **clients**  basic client component that use the gRPC protocol
* **services** basic service implementations for connecting via the gRPC

