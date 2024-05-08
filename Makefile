# You can build all images
build-all:
	./piglet-bills/scripts/build_images.sh

# You can run the app
run-app:
	docker-compose -f docker-compose.yml up --build

