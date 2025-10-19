# Use your existing docker-compose setup
FROM docker/compose:latest

# Copy your entire project
COPY . /app
WORKDIR /app

# Expose ports
EXPOSE 3000 8080

# Start everything with docker-compose
CMD ["docker-compose", "up", "--build"]
