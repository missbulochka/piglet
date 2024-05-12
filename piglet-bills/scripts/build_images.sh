#!/usr/bin/env bash

# You can build docker-image base (*.build.Dockerfile). This image is common for all services.
docker build \
        -t piglet-bills_base:0.1.0 \
        -f piglet-bills/.devcontainer/base.build.Dockerfile \
        ./

# You can build docker-image (*.Dockerfile)
docker build \
        -t piglet-bills:0.1.0 \
        -f piglet-bills/docker/piglet-bills.Dockerfile \
        ./