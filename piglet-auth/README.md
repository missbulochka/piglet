# Piglet authorization service

This is the README file for the Piglet Authorization Service,
a service that handles authorization in the Piglet finance
management application.

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
        -t piglet-auth_base:0.1.0 \
        -f piglet-auth/.devcontainer/piglet-auth.build.Dockerfile \
        ./

# start the dev container
devcontainer up --workspace-folder piglet-auth
```

You can run the application with the following command:
```bash
devcontainer exec --workspace-folder piglet-auth go run /workspaces/dev_piglet/cmd/main.go
```
