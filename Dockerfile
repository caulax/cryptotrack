FROM golang:1.23-alpine3.20 AS backend

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /build/cryptotrack


FROM alpine:3.20.3

WORKDIR /srv

COPY htmx /srv/htmx/
COPY static /srv/static/

COPY --from=backend /build/cryptotrack /srv/cryptotrack

ENTRYPOINT ["/srv/cryptotrack"]
