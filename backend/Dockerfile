# Simplified Go backend
FROM --platform=linux/amd64 golang:1.21-alpine

# Set working directory
WORKDIR /app

# Install minimal dependencies
RUN apk add --no-cache sqlite gcc musl-dev

# Copy everything
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o servers-listing main.go

# Expose port
EXPOSE 8080

# Run the application
CMD ["./servers-listing"]
