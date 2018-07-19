FROM golang:alpine as build

RUN apk update && \
    apk add make git ca-certificates

ADD . /go/src/github.com/subnova/vault-sidekick
RUN cd /go/src/github.com/subnova/vault-sidekick && make static test

FROM alpine:3.5
MAINTAINER Dale Peakall <dpeakall@thoughtworks.com>

RUN apk update && \
    apk add ca-certificates bash

RUN adduser -D vault

COPY --from=build /go/src/github.com/subnova/vault-sidekick/bin/vault-sidekick /vault-sidekick

USER vault

ENTRYPOINT [ "/vault-sidekick" ]
