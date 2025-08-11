# Multi-stage build for Chrome uTLS Template Generator
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chrome-utls-gen .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S chromeutls && \
    adduser -u 1001 -S chromeutls -G chromeutls

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/chrome-utls-gen .

# Create templates directory
RUN mkdir -p /app/templates && \
    chown -R chromeutls:chromeutls /app

# Switch to non-root user
USER chromeutls

# Expose port (if needed for web interface)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./chrome-utls-gen"]

# Default command
CMD ["--help"]
