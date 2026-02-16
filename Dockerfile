# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy dependency files first (for better layer caching)
COPY go.mod go.sum ./

# Download dependencies (cached layer - only re-runs if go.mod/go.sum changes)
RUN go mod download && go mod verify

# Copy source code (cached layer invalidated when code changes)
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o ./bin/main \
    ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o ./bin/migrate \
    ./cmd/migration

# Final stage
FROM alpine:3.23 AS final

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/bin/main main
COPY --from=builder /app/bin/migrate migrate

# Non-root user for security
RUN adduser -D -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

ENV GIN_MODE=release
EXPOSE 8888

CMD ["./main"]
