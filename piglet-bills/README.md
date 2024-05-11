# Piglet-bills service

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
        -f .devcontainer/base.build.Dockerfile \
        ./

# start the dev container
devcontainer up --workspace-folder .
```

You can run the application with the following command:
```bash
devcontainer exec --workspace-folder . go run /workspaces/dev_piglet/cmd/main.go
```

To stop containers use:
```bash
docker stop dev_piglet-bills bills-psql && docker rm dev_piglet-bills bills-psql
```

## Database (PostgreSQL)

To manage the databases for the Piglet-bill service locally, you can use the following commands
(you should export or fill POSTGRES_USER, POSTGRES_PASSWORD, HOST, PORT, POSTGRES_DB).

### Work with database

Use the following commands to create and delete database required for the Piglet-bills service:
```bash
# create database
docker exec \
        -it bills_psql \
        createdb --username=$POSTGRES_USER --owner=$POSTGRES_USER \
        Accounting

# delete database 
docker exec \
        -it bills_psql \
        dropdb --username=$POSTGRES_USER \
        Accounting
```

### Work with migrations

To run migrations for the Piglet-bills service locally, use the following commands.

```bash
# migrate up
migrate -path piglet-bills/migration \
        -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$HOST:$PORT/$POSTGRES_DB?sslmode=disable" \
        -verbose \
        up

# migrate down
migrate -path piglet-bills/migration \
        -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$HOST:$PORT/$POSTGRES_DB?sslmode=disable" \
        -verbose \
        down
```