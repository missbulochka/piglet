# Generate Client and Server code using proto file
generate-trans-pb:
	./scripts/gen_grpc.sh

# Up devcontainer
up-dev-trans:
	devcontainer up --workspace-folder .

# Down devcontainer (and dependencies)
rm-dev-trans:
	docker stop \
		transactions-psql \
		dev_piglet-transactions \
		bills-psql \
		piglet-bills && \
	docker rm \
		transactions-psql \
        dev_piglet-transactions \
        bills-psql \
        piglet-bills

# Run devcontainer
run-dev-trans:
	devcontainer exec --workspace-folder . go run /workspaces/dev_piglet/cmd/main.go

# Migrate up
up-migrate-trans:
	migrate -path migration \
		-database "postgresql://postgres:pass1234@transactions-psql:5432/Accounting?sslmode=disable" \
		-verbose \
		up

# Migrate down
down-migrate-trans:
	migrate -path migration \
	-database "postgresql://postgres:pass1234@transactions-psql:5432/Accounting?sslmode=disable" \
	-verbose \
	down
