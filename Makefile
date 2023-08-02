BUILD_VERSION ?= $(shell git describe --tags --dirty)
DOCKER_IMAGE  ?= ubercontroller
DOCKER_TAG    ?= develop
LDFLAGS       ?=
ACR_REPO      ?= odysseyprod.azurecr.io

all: build

vendor: go.mod
	go mod vendor
gen: vendor
	go generate ./...

gen-clean:
	find . -type f \( -name "*.mus.go" -o -name "*.autogen.go" \) | xargs rm

build: gen
	go build -ldflags "${LDFLAGS} -X main.version=${BUILD_VERSION}" -buildvcs=false -trimpath -o ./bin/ubercontroller ./cmd/service
	cd plugins && make

run: build
	./bin/ubercontroller

test:
	mkdir -p build
	go test -v -race -coverprofile=build/coverage.txt $$(go list ./... | grep -v  -E "ubercontroller/(build|cmd|docs)")

build-docs:
	go run github.com/swaggo/swag/cmd/swag@v1.8.8 init -g api.go -d universe/node,./,universe/streamchat -o build/docs/

docker-build: DOCKER_BUILDKIT=1
docker-build:
	docker build --build-arg BUILD_VERSION=${BUILD_VERSION} -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

docker-push-acr:
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${ACR_REPO}/${DOCKER_IMAGE}:${DOCKER_TAG}
	# az acr login -n odysseyprod
	docker push ${ACR_REPO}/${DOCKER_IMAGE}:${DOCKER_TAG}

# docker run ...
docker: docker-build
	docker run --rm ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: build run test docker docker-build
