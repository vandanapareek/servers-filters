# Multi-stage build for Go backend
FROM --platform=linux/amd64 golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git, ca-certificates, and gcc (needed for go mod download and CGO)
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o servers-listing main.go

# Final stage
FROM --platform=linux/amd64 alpine:latest

# Install ca-certificates and sqlite
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/servers-listing .

# Copy the database file
COPY data/servers.db ./data/servers.db

# Create data directory if it doesn't exist
RUN mkdir -p data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./servers-listing"]
