all: build

build:
	go build -o ./bin/controller ./cmd/service

run: build
	./bin/controller

test:
	go test -v -race ./...

# docker run ...
docker:
	echo "Not implemented"

.PHONY: build run test docker
