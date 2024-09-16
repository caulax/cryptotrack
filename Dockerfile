FROM golang:1.23-alpine3.20 as backend

ADD . /build
WORKDIR /build

RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1

RUN go build -o /build/cryptotrack


FROM alpine:3.20.3

WORKDIR /srv

COPY htmx /srv/htmx/
COPY static /srv/static/

COPY --from=backend /build/cryptotrack /srv/cryptotrack

ENTRYPOINT ["/srv/cryptotrack"]
