generate_proto:
	@protoc -I=btrace-idl/proto --go_out=proto-gen .\btrace-idl\proto\v2\proxy.proto
	@protoc -I=btrace-idl/proto --go-grpc_out=proto-gen .\btrace-idl\proto\v2\proxy.proto

	@protoc -I=btrace-idl/proto --go_out=proto-gen .\btrace-idl\proto\v2\storage.proto
	@protoc -I=btrace-idl/proto --go-grpc_out=proto-gen .\btrace-idl\proto\v2\storage.proto
