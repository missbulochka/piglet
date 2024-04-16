# Generate Client and Server code using proto file
generate:
	mkdir -p internal/proto/grpc/auth/service
	protoc -I=. \
		--go_out=. \
        --go-grpc_out=. \
        api/proto/grpc/auth/auth.proto
