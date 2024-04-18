#!/usr/bin/env bash

# Generate Client and Server code for auth-service using proto file
mkdir -p piglet-auth/api/proto/gen
protoc -I=. \
    --go_out=. \
    --go-grpc_out=. \
    piglet-auth/api/proto/auth.proto
