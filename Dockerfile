# syntax=docker/dockerfile:1.4
FROM golang:1.19-alpine3.16 as build

RUN apk add --update --no-cache gcc binutils-gold musl-dev

WORKDIR /project

# Seperate step to allow docker layer caching
COPY go.* ./
RUN go mod download

COPY . ./


# extra ldflag to make sure it works with alpine/musl

RUN go generate ./... && go build -trimpath -ldflags "-extldflags '-fuse-ld=bfd'" -o ./bin/ubercontroller ./cmd/service
#RUN go build -o ./bin/ubercontroller ./cmd/service


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
