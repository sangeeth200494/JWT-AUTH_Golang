# Multi-stage production Dockerfile
# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy application source code
COPY . .

# Build standalone static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o jwt-auth-server main.go

# Stage 2: Minimal hardened runtime image
FROM alpine:3.20

# Create non-root user for principle of least privilege
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy system certs, application binary, and frontend assets
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/jwt-auth-server /app/jwt-auth-server
COPY --from=builder /app/static /app/static
COPY --from=builder /app/.env.example /app/.env.example

# Set user permissions
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 4000

# Container healthcheck
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:4000/Home || exit 1

ENTRYPOINT ["/app/jwt-auth-server"]
