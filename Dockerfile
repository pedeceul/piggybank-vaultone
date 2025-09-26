# syntax=docker/dockerfile:1.6

FROM golang:1.24 AS builder
WORKDIR /src

# Leverage Docker layer caching
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /out/api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /out/api /bin/api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/bin/api"]

