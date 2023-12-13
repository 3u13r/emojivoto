IMAGE_TAG ?= coco-1

.PHONY: package protoc test

target_dir := target

clean:
	rm -rf $(target_dir)
	mkdir -p $(target_dir)

protoc:
	cd tools && DOCKER_BUILDKIT=1 docker build -o .. -f Dockerfile.gen-proto ..

package: compile build-container

build-container:
	docker build .. -t "ghcr.io/3u13r/$(svc_name):$(IMAGE_TAG)" --build-arg svc_name=$(svc_name)

compile:
	GOOS=linux go build -v -o $(target_dir)/$(svc_name) cmd/server.go

test:
	go test ./...

run:
	go run cmd/server.go
