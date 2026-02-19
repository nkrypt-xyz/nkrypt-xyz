FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/nkrypt-server \
    ./cmd/server

FROM alpine:3.19

RUN addgroup -S nkrypt && adduser -S nkrypt -G nkrypt
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/nkrypt-server .
COPY --from=builder /build/migrations ./migrations

RUN chown -R nkrypt:nkrypt /app
USER nkrypt

EXPOSE 9041

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9041/healthz || exit 1

ENTRYPOINT ["./nkrypt-server"]

