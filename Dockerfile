# Step 1: Modules caching and build stage
FROM golang:1.26.4-alpine AS builder

# Install certs and tzdata for security and timezone handling
RUN apk add --no-cache ca-certificates tzdata

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN swag init --parseDependency --parseInternal

# Build the Go application
# CGO_ENABLED=0 creates a static binary
# -ldflags="-s -w" strip debugging information and symbols to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o refinery main.go

# Step 2: Minimal runner stage
FROM alpine:3.20 AS runner

# Add a non-root user for security best practices
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy SSL certificates and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy compiled binary and web UI resources from builder with non-root ownership
COPY --chown=appuser:appgroup --from=builder /app/refinery /app/refinery
COPY --chown=appuser:appgroup --from=builder /app/internal/web/views /app/internal/web/views
COPY --chown=appuser:appgroup --from=builder /app/internal/web/public /app/internal/web/public

# Use the non-root user
USER appuser

# Expose port 7003 (default port from example.env)
EXPOSE 7003

# Run the server command by default
ENTRYPOINT ["/app/refinery"]
CMD ["serve"]
