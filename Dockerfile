# syntax=docker/dockerfile:1.4
ARG BUILD_VERSION
FROM golang:1.20-alpine3.17 as build
ARG BUILD_VERSION

RUN apk add --update --no-cache gcc make binutils-gold musl-dev

WORKDIR /project

# Seperate step to allow docker layer caching
COPY go.* ./
RUN go mod download

COPY . ./

# extra ldflag to make sure it works with alpine/musl
ENV LDFLAGS="-extldflags '-fuse-ld=bfd'" BUILD_VERSION=${BUILD_VERSION}
RUN make build

# Runtime image
FROM alpine:3.16 as runtime

LABEL org.opencontainers.image.source=https://github.com/momentum-xyz/ubercontroller
LABEL org.opencontainers.image.description="Controller of Odyssey Momentum"
LABEL org.opencontainers.image.licenses=AGPL-3.0-only

COPY --link ./seed/data /srv/seed/data

COPY --from=build /project/bin/ubercontroller /srv/ubercontroller

WORKDIR /srv
CMD ["/srv/ubercontroller"]
