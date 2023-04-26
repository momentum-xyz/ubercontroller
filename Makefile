DOCKER_IMAGE="ubercontroller"
DOCKER_TAG="develop"

all: build

vendor: go.mod
	go mod vendor
gen: vendor
	go generate ./...

gen-clean:
	find . -type f \( -name "*.mus.go" -o -name "*.autogen.go" \) | xargs rm

build: gen
	go build -trimpath -o ./bin/ubercontroller ./cmd/service
	cd plugins && make

run: build
	./bin/ubercontroller

test:
	go test -v -race ./...

build-docs:
	swag init -g api.go -d universe/node,./,universe/streamchat -o build/docs/ --parseDependency --parseInternal

docker-build: DOCKER_BUILDKIT=1
docker-build:
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

# docker run ...
docker: docker-build
	docker run --rm ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: build run test docker docker-build
