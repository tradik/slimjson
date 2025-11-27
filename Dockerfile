# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o slimjson ./cmd/slimjson

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/slimjson .

# Create a non-root user
RUN addgroup -g 1000 slimjson && \
    adduser -D -u 1000 -G slimjson slimjson && \
    chown -R slimjson:slimjson /app

USER slimjson

# Expose port for HTTP daemon
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default to daemon mode
ENTRYPOINT ["/app/slimjson"]
CMD ["-d", "-port", "8080"]
