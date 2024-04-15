# Minimum viable product

## Prerequisites

### devcontainers/cli

```bash
npm install -g @devcontainers/cli
```

## Requirements

- go 1.21.3
- docker 26.0.0
- docker-compose 1.29.2 
- devcontainers/cli (optional)

## Build

### Dev Container Cli

You can build a project with [devcontainers](https://containers.dev/) in an easy and convenient way.
Your IDE or code editor can run and attach to devcontainer.

You can use devcontainers/cli to set up environment and build the project manually via bash:
```
docker build \
    -t hello-world_build-base:0.1.0 \
    -f ./.devcontainer/build.Dockerfile \
    ./
    
devcontainer up --workspace-folder .

devcontainer exec --workspace-folder . \
    go run /workspaces/dev_hello-world/cmd/main.go
```

## Run
You can run the application with docker (call from directory `hello_world`):
```
docker build \
    -t hello-world_build-base:0.1.0 \
    -f ./.devcontainer/build.Dockerfile \
    ./

docker build \
    -t app_hello-world:0.1.0 \
    -f ./docker/Dockerfile \
    ./

docker-compose -f ./deploy/docker-compose.yml up
```