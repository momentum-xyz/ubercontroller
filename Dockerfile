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


RUN apk add --update --no-cache python3 make g++
#temporary, add nodejs and polkadot package

RUN apk add --update --no-cache nodejs npm 
COPY ./nodejs /srv/nodejs
WORKDIR /srv/nodejs/check-nft
RUN npm i

COPY --from=build /project/bin/ubercontroller /srv/ubercontroller
COPY --link ./seed/data /srv/seed/data

WORKDIR /srv
CMD ["/srv/ubercontroller"]
