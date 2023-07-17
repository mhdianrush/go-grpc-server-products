GO_WORKSPACE := $(GOPATH)/src/03-projects

protoc:
	protoc --proto_path=protos protos/*.proto --go_out=$(GO_WORKSPACE) --go-grpc_out=$(GO_WORKSPACE) --experimental_allow_proto3_optional
	@echo "Protoc Compile Successfully"