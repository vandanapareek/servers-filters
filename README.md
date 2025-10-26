# Servers Filters

A full-stack web application for filtering and managing server listings with advanced search capabilities.

## Project Overview

This project provides a RESTful API backend built with Go and a modern Vue.js frontend for filtering server data. Users can search and filter servers based on various criteria including RAM, storage, location, and price.

## Features

- **Server Filtering**: Filter by RAM, storage, location, and hard disk type
- **Pagination**: Efficient data pagination for large datasets
- **Real-time Search**: Fast server search and filtering
- **Responsive UI**: Modern Vue.js frontend with responsive design
- **RESTful API**: Clean Go backend with comprehensive API endpoints

## Tech Stack

### Backend
- **Go** with Chi router
- **SQLite** database
- **Docker** containerization

### Frontend
- **Vue.js 3** with Composition API
- **Vite** build tool
- **Axios** for API calls
- **Nginx** web server

## Local Development Setup

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for backend development)
- Node.js 20+ (for frontend development)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd servers-filters
   ```

2. **Start with Docker Compose**
   ```bash
   docker-compose up
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8081

## Port Requirements & Troubleshooting

### Required Ports
This application requires the following ports to be available:
- **Port 3000**: Frontend (Vue.js app)
- **Port 8081**: Backend API

### Port Conflicts
If you encounter port conflicts, you have several options:

#### Option 1: Stop conflicting services
```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8081

# Stop the conflicting processes
kill -9 <PID>
```

#### Option 2: Change ports in docker-compose.yml
```yaml
services:
  frontend:
    ports:
      - "3001:80"  # Use port 3001 instead
  backend:
    ports:
      - "8082:8080"  # Use port 8082 instead
```

### Verification
After starting, verify all services are running:
```bash
docker-compose ps
curl http://localhost:3000  # Frontend
curl http://localhost:8081/servers  # Backend
```

### Manual Development Setup

#### Backend Setup
```bash
cd backend
go mod download
go run main.go
```

#### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## API Documentation

### Postman Collection
Import the provided Postman collection to test all API endpoints:

**File**: [`docs/ServersFilters.postman_collection.json`](./docs/ServersFilters.postman_collection.json)

**Setup**:
1. Open Postman
2. Import the collection file
3. Set base URL to your backend URL
4. Test all endpoints with examples

### Swagger/OpenAPI Documentation
Interactive API documentation with Swagger UI:

**File**: [`docs/swagger.yaml`](./docs/swagger.yaml)

**Setup**:
1. Go to [Swagger Editor](https://editor.swagger.io/)
2. Copy the content from [`docs/swagger.yaml`](./docs/swagger.yaml) file
3. Paste it into the Swagger Editor
4. Test all endpoints interactively with full documentation


## AWS Deployment

The application is deployed on AWS ECS with the following URLs:

- **Frontend**: http://54.237.24.242/
- **Backend API**: http://3.236.4.27:8080/servers

### Deployment Process

1. **Build and push images**
   ```bash
   ./simple-aws-deploy.sh
   ```

2. **Update ECS task definitions** in AWS Console:
   - Go to ECS → Task Definitions
   - Update backend task definition with: `916591320534.dkr.ecr.us-east-1.amazonaws.com/servers-filters-backend:arm64`
   - Update frontend task definition with: `916591320534.dkr.ecr.us-east-1.amazonaws.com/servers-filters-frontend:arm64`

3. **Deploy via AWS Console**:
   - Go to ECS → Services
   - Update service to use new task definition revision

### Infrastructure
- **ECS** for container orchestration
- **ECR** for Docker image storage

### Environment Variables
- Frontend: `VITE_API_BASE_URL` (backend API URL)
- Backend: Database and port configuration


## Testing

Run backend tests:
```bash
cd backend
go test ./...
```
