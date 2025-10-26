FROM --platform=linux/amd64 golang:1.21-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache sqlite gcc musl-dev

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o servers-filters main.go

# Expose port
EXPOSE 8080

# Run the application
CMD ["./servers-filters"]
