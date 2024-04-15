# Call from directory `hello_world`

# Generate Client and Server code using proto file
generate:
	mkdir -p pkg/user_v1
	protoc --go_out=./ --go_opt=paths=import \
    		--go-grpc_out=./ --go-grpc_opt=paths=import \
            api/proto/service.proto

# You can build a project with devcontainers in an easy and convenient way.
app-build:
	docker build \
        -t hello-world_build-base:0.1.0 \
        -f ./.devcontainer/build.Dockerfile \
        ./
	devcontainer up --workspace-folder .
	devcontainer exec --workspace-folder . \
        go run /workspaces/dev_hello-world/cmd/main.go

# You can run the application with docker
run-app:
	docker build \
        -t hello-world_build-base:0.1.0 \
        -f ./.devcontainer/build.Dockerfile \
        ./

	docker-compose -f ./deploy/docker-compose.yml up --build

# You can stop docker-compose with
down-app:
	docker-compose -f ./deploy/docker-compose.yml down
