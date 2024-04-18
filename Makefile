# Generate Client and Server code using proto file
generate-proto:
	mkdir -p internal/proto/grpc/auth/service
	protoc -I=. \
		--go_out=. \
        --go-grpc_out=. \
        api/proto/grpc/auth/auth.proto

# You can build all *.build.Dockerfile
build-all-build:
	docker build \
        -t piglet-auth_base:0.1.0 \
        -f ./devcontainers/piglet-auth/.devcontainer/piglet-auth.build.Dockerfile \
        ./

# You can build all *.Dockerfile
build-all-srv:
	docker build \
		-t piglet-auth:0.1.0 \
		-f ./docker/piglet-auth.Dockerfile \
		./

# You can build an auth-service with devcontainers in an easy and convenient way.
dev-auth-service:
	devcontainer up --workspace-folder devcontainers/piglet-auth
	devcontainer exec --workspace-folder devcontainers/piglet-auth \
            go run /workspaces/dev_piglet/cmd/piglet-auth/main.go
