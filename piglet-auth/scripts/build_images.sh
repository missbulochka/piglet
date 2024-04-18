#!/usr/bin/env bash

# You can build docker-image base (*.build.Dockerfile)
docker build \
        -t piglet-auth_base:0.1.0 \
        -f piglet-auth/.devcontainer/piglet-auth.build.Dockerfile \
        ./

# You can build docker-image (*.Dockerfile)
docker build \
        -t piglet-auth:0.1.0 \
        -f piglet-auth/docker/piglet-auth.Dockerfile \
        ./
