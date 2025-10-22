# Multi-stage build for Vue.js frontend
FROM --platform=linux/amd64 node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage with nginx
FROM --platform=linux/amd64 nginx:alpine

# Accept backend URL at build time to bake per-environment config
ARG BACKEND_URL
ENV BACKEND_URL=${BACKEND_URL}

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Inject BACKEND_URL into nginx.conf
RUN sh -c 'sed -i "s|\${BACKEND_URL}|${BACKEND_URL}|g" /etc/nginx/nginx.conf'

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80 || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
