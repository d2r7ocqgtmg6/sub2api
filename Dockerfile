# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
# Note: removed GOARCH=amd64 hardcode so cross-compilation works on arm64 hosts too
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -o sub2api \
    ./main.go

# Final stage — minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the compiled binary
COPY --from=builder /app/sub2api /sub2api

# Expose the default port
# Changed from 8080 to 8088 to avoid conflicts with other local services
EXPOSE 8088

# Run as non-root user for security
USER 65534:65534

ENTRYPOINT ["/sub2api"]
