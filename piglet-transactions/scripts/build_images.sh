#!/usr/bin/env bash

# You can build docker-image base (*.build.Dockerfile). This image is common for all services.
docker build \
        -t piglet-transactions_base:0.1.0 \
        -f piglet-transactions/.devcontainer/base.build.Dockerfile \
        ./

# You can build docker-image (*.Dockerfile)
docker build \
        -t piglet-transactions:0.1.0 \
        -f piglet-transactions/docker/piglet-transactions.Dockerfile \
        ./