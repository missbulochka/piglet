# Piglet authorization service

This README file is for the piglet-bills service, which manages bills and goals.

## Prerequisites

### devcontainers/cli
Install the Dev Container CLI:

```bash
npm install -g @devcontainers/cli
```

## Build

### Dev Container Cli

You can build the project with [devcontainers](https://containers.dev/) in an easy and convenient way.
Your IDE or code editor can run and attach to devcontainer.

You can use devcontainers/cli to set up environment and build the project manually via bash:
```bash
# build the docker image
docker build \
        -t piglet_base:0.1.0 \
        -f piglet-bills/.devcontainer/base.build.Dockerfile \
        ./

# start the dev container
devcontainer up --workspace-folder piglet-bills
```

You can run the application with the following command:
```bash
devcontainer exec --workspace-folder piglet-bills go run /workspaces/dev_piglet/cmd/main.go
```