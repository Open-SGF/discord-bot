FROM golang:1.25

RUN mkdir -p /cache/go-cache /cache/go-modcache && \
    chmod -R 777 /cache

ENV GOCACHE=/cache/go-cache
ENV GOMODCACHE=/cache/go-modcache

WORKDIR /asset-input

ENV CGO_ENABLED=0
