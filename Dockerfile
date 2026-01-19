# Build Stage
# Using 1.25-alpine as requested/implied by go.mod, assuming it exists in 2026. 
# If not, 'latest' or specific version available at the time should be used.
FROM golang:1.25-alpine AS builder

# Install build dependencies
# gcc and musl-dev are required for go-sqlite3 (CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for sqlite3
RUN CGO_ENABLED=1 go build -o icomment .

# Binary Stage
# Used for exporting the binary to host
FROM scratch AS binary
COPY --from=builder /app/icomment /

# Runtime Stage
FROM alpine:latest AS runtime

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/icomment .

# Create volume directory for database
RUN mkdir -p /data
VOLUME /data

# Expose ports
# 7001: Public API
# 7002: Admin Panel
EXPOSE 7001 7002

# Run the application
# Default database path set via ENV for easier docker-compose configuration
ENV ICOMMENT_DB=/data/comments.db
ENTRYPOINT ["/app/icomment"]
