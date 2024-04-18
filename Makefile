# Generate Client and Server code using proto file
generate-pb-go:
	./piglet-auth/scripts/gen_grpc.sh

# You can build all images
build-all:
	./piglet-auth/scripts/build_images.sh
