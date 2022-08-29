all: build

build:
	go build -o ./bin/ubercontroller ./cmd/service

run: build
	./bin/ubercontroller

test:
	go test -v -race ./...

# docker run ...
docker:
	echo "Not implemented"

.PHONY: build run test docker
