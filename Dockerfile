FROM alpine:3.5
MAINTAINER Dale Peakall <dpeakall@thoughtworks.com>

RUN apk update && \
    apk add ca-certificates bash

RUN adduser -D vault

ADD bin/vault-sidekick /vault-sidekick

USER vault

ENTRYPOINT [ "/vault-sidekick" ]
