# Piglet-transactions service

This README file is for the piglet-transactions service, which manages all transactions.

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
        -t piglet-bills_base:0.1.0 \
        -f piglet-bills/.devcontainer/base.build.Dockerfile \
        ./

# You can build docker-image (*.Dockerfile)
docker build \
        -t piglet-bills:0.1.0 \
        -f piglet-bills/docker/piglet-bills.Dockerfile \
        ./

docker build \
        -t piglet-transactions_base:0.1.0 \
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
docker stop \
		transactions-psql \
		dev_piglet-transactions \
		bills-psql \
		piglet-bills
```

To remove containers use:
```bash
docker rm \
		transactions-psql \
        dev_piglet-transactions \
        bills-psql \
        piglet-bills
```

## Database (PostgreSQL)

To manage the databases for the Piglet-transactions service locally, you can use the following commands
(you should export or fill POSTGRES_USER, POSTGRES_PASSWORD, HOST, PORT, POSTGRES_DB).

### Work with database

Use the following commands to create and delete database required for the Piglet-transactions service:
```bash
# create database
docker exec \
        -it transactions-psql \
        createdb --username=$POSTGRES_USER --owner=$POSTGRES_USER \
        Accounting

# delete database 
docker exec \
        -it transactions-psql \
        dropdb --username=$POSTGRES_USER \
        Accounting
```

### Work with migrations

To run migrations for the Piglet-transactions service locally, use the following commands.

```bash
# migrate up
migrate -path piglet-transactions/migration \
        -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$HOST:$PORT/$POSTGRES_DB?sslmode=disable" \
        -verbose \
        up

# migrate down
migrate -path piglet-transactions/migration \
        -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$HOST:$PORT/$POSTGRES_DB?sslmode=disable" \
        -verbose \
        down
```