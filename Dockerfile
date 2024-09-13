FROM golang:1.23-alpine as backend

ADD . /build
WORKDIR /build

RUN go build -o /build/cryptotrack


FROM scratch

COPY --from=backend /build/cryptotrack /srv/cryptotrack

WORKDIR /srv
ENTRYPOINT ["/srv/cryptotrack"]
