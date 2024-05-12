#!/usr/bin/env bash

# Generate Client and Server code for bills-service using proto file
mkdir -p api/proto/gen
protoc -I=. \
    --go_out=. \
    --go-grpc_out=. \
    api/proto/accounting.proto