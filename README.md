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
- **Redis** caching (optional)
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
   - Redis: localhost:6380

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

## API Endpoints

- `GET /servers` - List servers with filtering
- `GET /servers/{id}` - Get server by ID
- `GET /locations` - Get available locations
- `GET /metrics` - Get server statistics

## AWS Deployment

The application is deployed on AWS ECS with the following URLs:

- **Frontend**: http://44.220.245.64/
- **Backend API**: http://3.237.50.32:8080/servers

### Deployment Process

1. **Build and push images**
   ```bash
   ./simple-aws-deploy.sh
   ```

2. **Update ECS task definitions** with new image URIs
3. **Deploy to ECS services**


## Testing

Run backend tests:
```bash
cd backend
go test ./...
```
