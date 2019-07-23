.PHONY: all build clean install uninstall fmt simplify check run test

install:
	@go install main.go

run: install
	@go run main.go

test:
	@go test ./test/...

protogen:
	protoc --go_out=plugins=grpc:. protos/commandable.proto