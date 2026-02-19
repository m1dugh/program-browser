FROM golang:1.25-alpine3.23 AS builder

WORKDIR /build

RUN --mount=type=bind,src=go.mod,dst=go.mod \
    --mount=type=bind,src=go.sum,dst=go.sum \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN --mount=type=bind,src=.,dst=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /program-browser ./cmd/program-browser

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /program-browser /app

ENTRYPOINT ["/app/program-browser"]
