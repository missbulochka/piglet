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
		-database "postgresql://postgres:pass1234@127.0.0.1:5432/Transactions?sslmode=disable" \
		-verbose \
		up

# Migrate down
down-migrate-trans:
	migrate -path migration \
	-database "postgresql://postgres:pass1234@127.0.0.1:5432/Transactions?sslmode=disable" \
	-verbose \
	down
