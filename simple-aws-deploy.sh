#!/bin/bash

# Simple AWS ECR deployment for existing ECS setup
# This script just builds and pushes images, then gives you the URIs

set -e

# Configuration
AWS_REGION="us-east-1"
AWS_ACCOUNT_ID="916591320534"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 Building images for your existing ECS setup${NC}"

# Login to ECR
echo -e "${YELLOW}🔐 Logging in to ECR...${NC}"
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

# Build and push backend (no changes needed)
echo -e "${YELLOW}🏗️  Building backend image for AMD64...${NC}"
cd backend
docker build -t servers-filters-backend:amd64 .
docker tag servers-filters-backend:amd64 ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/servers-filters-backend:amd64
docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/servers-filters-backend:amd64
cd ..

# Build and push frontend (using existing Dockerfile)
echo -e "${YELLOW}🏗️  Building frontend image for AMD64...${NC}"
cd frontend
docker build -t servers-filters-frontend:amd64 .
docker tag servers-filters-frontend:amd64 ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/servers-filters-frontend:amd64
docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/servers-filters-frontend:amd64
cd ..

echo -e "${GREEN}🎉 Images built and pushed successfully!${NC}"
echo -e "${GREEN}📋 Your ECR Image URIs:${NC}"
echo -e "Backend:  ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/servers-filters-backend:amd64"
echo -e "Frontend: ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/servers-filters-frontend:amd64"
echo ""
echo -e "${YELLOW}💡 Next steps:${NC}"
echo -e "1. Copy the image URIs above"
echo -e "2. Update your ECS task definitions with these URIs"
echo -e "3. For frontend task, add environment variable: VITE_API_BASE_URL=https://your-backend-alb-url"
echo -e "4. Update your ECS services to use the new task definitions"
