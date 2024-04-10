# Generate Client and Server code using proto file
# Call from directory `hello_world`
PHONY: generate
generate:
	mkdir -p pkg/user_v1
	protoc --go_out=./ --go_opt=paths=import \
    		--go-grpc_out=./ --go-grpc_opt=paths=import \
            api/proto/service.proto
