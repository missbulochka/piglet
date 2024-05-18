# You can build docker-image (*.Dockerfile)
docker build \
        -t piglet-gateway:0.1.0 \
        -f piglet-gateway/docker/piglet-gateway.Dockerfile \
        ./