# Multi-stage build for tf-iamgen
# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && \
    go mod verify

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME=""
ARG GIT_COMMIT=""

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-w -s \
              -X 'main.Version=${VERSION}' \
              -X 'main.BuildTime=${BUILD_TIME}' \
              -X 'main.GitCommit=${GIT_COMMIT}'" \
    -o tf-iamgen .

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S tf-iamgen && adduser -S tf-iamgen -G tf-iamgen

WORKDIR /home/tf-iamgen

# Copy binary from builder
COPY --from=builder /app/tf-iamgen /usr/local/bin/tf-iamgen

# Change ownership
RUN chown tf-iamgen:tf-iamgen /usr/local/bin/tf-iamgen

# Switch to non-root user
USER tf-iamgen

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD tf-iamgen version || exit 1

# Set entrypoint
ENTRYPOINT ["tf-iamgen"]
CMD ["help"]
