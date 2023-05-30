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
	mkdir -p build
	go test -v -race -coverprofile=build/coverage.txt $$(go list ./... | grep -v  -E "ubercontroller/(build|cmd|docs)")

build-docs:
	swag init -g api.go -d universe/node,./,universe/streamchat -o build/docs/

docker-build: DOCKER_BUILDKIT=1
docker-build:
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

# docker run ...
docker: docker-build
	docker run --rm ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: build run test docker docker-build
