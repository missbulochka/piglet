# Up devcontainer
up-dev-bills:
	devcontainer up --workspace-folder .

# Down devcontainer (and dependencies)
rm-dev-bills:
	docker stop dev_piglet-bills bills-psql && docker rm dev_piglet-bills bills-psql

# Run devcontainer
run-dev-bills:
	devcontainer exec --workspace-folder . go run /workspaces/dev_piglet/cmd/main.go

# Migrate up
up-migrate-bills:
	migrate -path migration \
            -database "postgresql://postgres:pass1234@bills-psql:5432/Accounting?sslmode=disable" \
            -verbose \
            up

# Migrate down
down-migrate-bills:
	migrate -path migration \
            -database "postgresql://postgres:pass1234@bills-psql:5432/Accounting?sslmode=disable" \
            -verbose \
            down
