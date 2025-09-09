# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.23.5

FROM golang:${GO_VERSION} AS builder
WORKDIR /workspace/services

# Copy the entire repo to leverage go.work and local modules
COPY . /workspace

# Enable reliable module caching
ENV CGO_ENABLED=0

# Build the authority-service server with required tags
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    go build -tags "postgres sqlc" -o /out/authority ./authority-service/cmd/server

# Build golang-migrate CLI with Postgres driver enabled for the current architecture
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /out/authority /app/authority
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

USER appuser
EXPOSE 8080

# GRPC_ADDR defaults to :8080; configure DB/IDP via envs
ENTRYPOINT ["/app/authority"]
