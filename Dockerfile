# syntax=docker/dockerfile:1.3
FROM golang:1.19-alpine3.16 as build

WORKDIR /project

# Seperate step to allow docker layer caching
COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o ./bin/ubercontroller ./cmd/service


# Runtime image
FROM alpine:3.16 as runtime

COPY --from=build /project/bin/ubercontroller /srv/ubercontroller

CMD ["/srv/ubercontroller"]
