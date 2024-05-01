# Generate all Client and Server code using proto file
generate-pb-go:
	./piglet-bills/scripts/gen_grpc.sh

# You can build all images
build-all:
	./piglet-bills/scripts/build_images.sh