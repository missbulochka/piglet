PHONY: generate
generate:
	mkdir -p pkg/user_v1
	protoc --go_out=pkg/user_v1 --go_opt=paths=import \
    		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=import \
            api/proto/service.proto
	mv pkg/user_v1/pkg/user_v1/* pkg/user_v1/
	rm -rf pkg/user_v1/pkg
