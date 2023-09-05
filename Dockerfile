# syntax=docker/dockerfile:1.4
ARG repo=ghcr.io/momentum-xyz/
ARG version_ui=develop
ARG BUILD_VERSION

###############
# Build stage #
###############
FROM golang:1.20-alpine3.17 as build
ARG BUILD_VERSION

RUN apk add --update --no-cache gcc make binutils-gold musl-dev

WORKDIR /project
ENV GOPATH /go
ENV GOCACHE /go-cache

# Seperate step to allow docker layer caching
COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY . ./

# extra ldflag to make sure it works with alpine/musl
ENV LDFLAGS="-s -w -extldflags '-fuse-ld=bfd'" BUILD_VERSION=${BUILD_VERSION}
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    make build


##################
# Runtime target #
##################
FROM alpine:3.16 as runtime

LABEL org.opencontainers.image.source=https://github.com/momentum-xyz/ubercontroller
LABEL org.opencontainers.image.description="Controller of Odyssey Momentum"
LABEL org.opencontainers.image.licenses=AGPL-3.0-only

COPY --link ./seed/data /srv/seed/data
COPY --link ./fonts /srv/fonts

COPY --from=build /project/bin/ubercontroller /srv/ubercontroller

WORKDIR /srv
EXPOSE 4000
CMD ["/srv/ubercontroller"]


##################
# Frontend image #
##################
FROM ${repo}ui-client:${version_ui} as ui-client


######################
# Embedded UI target #
######################
FROM runtime as monolith
env FRONTEND_SERVE_DIR=/srv/ui
COPY --from=ui-client --link /opt/srv /srv/ui


########################
# Default build target #
########################
FROM runtime as default
