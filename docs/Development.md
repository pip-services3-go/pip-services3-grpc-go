# Development and Testing Guide <br/> Pip.Services Google Remote Procedure Calls for Golang

This document provides high-level instructions on how to build and test the microservice.

* [Environment Setup](#setup)
* [Installing](#install)
* [Building](#build)
* [Testing](#test)
* [Release](#release)
* [Contributing](#contrib) 

## <a name="setup"></a> Environment Setup

This is a Golang project and you have to install Golang tools. 
You can download them from official Golang website: https://golang.org/dl/ 

After go is installed you can check it by running the following command:
```bash
go version
```

To work with GitHub code repository you need to install Git from: https://git-scm.com/downloads

## <a name="install"></a> Installing

After your environment is ready you can check out source code from the Github repository:
```bash
git clone git@github.com:pip-services/pip-services-grpc-go.git
```

Install proto-gen-go:
```bash
go install github.com/golang/protobuf/protoc-gen-go
```
Help for work with protobuf utility:

https://github.com/protocolbuffers/protobuf
https://github.com/golang/protobuf 

If you worked with the microservice before you can check out latest changes and update the dependencies:
# Update source code updates from github

```bash
go get -u
```

## <a name="build"></a> Building

The commons is written in Golang language.
So, if you make changes to the source code you need to compile it before running or committing to github.
The process will output compiled executable files.

```bash
go build
```

## <a name="test"></a> Testing

Before you execute tests you need to set configuration options in config.yaml file.
As a starting point you can use example from config.example.json:

```bash
copy config/config.example.yaml config/config.yaml
``` 

After that check all configuration options. Specifically, pay attention to connection options
for database and dependent microservices. For more information check [Configuration Guide](Configuration.md) 

Command to run unit tests and benchmarks as:
```bash
go test ./test/...
```

## <a name="release"></a> Release

Formal release process consistents of few steps. 
First of all it is required to tag guthub repository with a version number:

```bash
git commit -m "Emphasize our friendliness" testmod.go
git tag v1.0.1
git push --tags origin v1
```

Microservice releases additionally require generation and publishing 
binary packages at http://downloads.pipservices.org


## <a name="contrib"></a> Contributing

Developers interested in contributing should read the following instructions:

- [How to Contribute](http://www.pipservices.org/contribute/)
- [Guidelines](http://www.pipservices.org/contribute/guidelines)
- [Styleguide](http://www.pipservices.org/contribute/styleguide)
- [ChangeLog](../CHANGELOG.md)

> Please do **not** ask general questions in an issue. Issues are only to report bugs, request
  enhancements, or request new features. For general questions and discussions, use the
  [Contributors Forum](http://www.pipservices.org/forums/forum/contributors/).

It is important to note that for each release, the [ChangeLog](../CHANGELOG.md) is a resource that will
itemize all:

- Bug Fixes
- New Features
- Breaking Changes