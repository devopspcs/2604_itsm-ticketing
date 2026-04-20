#!/bin/bash
# Docker Deployment Script - ITSM Ticketing System
# Usage: ./deploy/docker-deploy.sh [dev|prod]

set -e

ENVIRONMENT=${1:-dev}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=========================================="
echo "ITSM Ticketing System - Docker Deployment"
echo "=========================================="
echo "Environment: $ENVIRONMENT"
echo "Project Dir: $PROJECT_DIR"
echo ""

# Validate environment
if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "prod" ]]; then
    echo "❌ Invalid environment. Use 'dev' or 'prod'"
    exit 1
fi

# Change to project directory
cd "$PROJECT_DIR"

# Step 1: Verify prerequisites
echo "📋 Step 1: Verifying prerequisites..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    exit 1
fi

echo "✅ Docker and Docker Compose found"
echo ""

# Step 2: Setup environment
echo "📋 Step 2: Setting up environment..."
if [ "$ENVIRONMENT" = "prod" ]; then
    if [ ! -f ".env.prod" ]; then
        echo "❌ .env.prod not found"
        echo "   Please create .env.prod with production configuration"
        exit 1
    fi
    ENV_FILE=".env.prod"
    COMPOSE_FILE="docker-compose.prod.yml"
    echo "✅ Using production environment"
else
    if [ ! -f ".env" ]; then
        echo "⚠️  .env not found, creating from .env.example"
        if [ -f ".env.example" ]; then
            cp .env.example .env
        else
            echo "❌ .env.example not found"
            exit 1
        fi
    fi
    ENV_FILE=".env"
    COMPOSE_FILE="docker-compose.yml"
    echo "✅ Using development environment"
fi
echo ""

# Step 3: Build Docker images
echo "📋 Step 3: Building Docker images..."
echo "   This may take a few minutes..."
if [ "$ENVIRONMENT" = "prod" ]; then
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build --no-cache
else
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build
fi
echo "✅ Docker images built successfully"
echo ""

# Step 4: Stop existing containers (if any)
echo "📋 Step 4: Stopping existing containers..."
if [ "$ENVIRONMENT" = "prod" ]; then
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down 2>/dev/null || true
else
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down 2>/dev/null || true
fi
echo "✅ Existing containers stopped"
echo ""

# Step 5: Start services
echo "📋 Step 5: Starting services..."
if [ "$ENVIRONMENT" = "prod" ]; then
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d
else
    docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d
fi
echo "✅ Services started"
echo ""

# Step 6: Wait for services to be ready
echo "📋 Step 6: Waiting for services to be ready..."
sleep 5

# Check PostgreSQL
echo "   Checking PostgreSQL..."
for i in {1..30}; do
    if docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres pg_isready -U itsm &>/dev/null; then
        echo "   ✅ PostgreSQL is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ❌ PostgreSQL failed to start"
        exit 1
    fi
    sleep 1
done

# Check Backend
echo "   Checking Backend..."
for i in {1..30}; do
    if curl -s http://localhost:8080/health &>/dev/null; then
        echo "   ✅ Backend is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ❌ Backend failed to start"
        exit 1
    fi
    sleep 1
done

# Check Frontend
echo "   Checking Frontend..."
for i in {1..30}; do
    if curl -s http://localhost:3000/health &>/dev/null; then
        echo "   ✅ Frontend is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ⚠️  Frontend health check timeout (may still be starting)"
        break
    fi
    sleep 1
done
echo ""

# Step 7: Display status
echo "📋 Step 7: Deployment Status"
echo "=========================================="
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps
echo "=========================================="
echo ""

# Step 8: Display access information
echo "📋 Access Information"
echo "=========================================="
echo "Frontend:  http://localhost:3000"
echo "Backend:   http://localhost:8080"
echo "API Docs:  http://localhost:8080/swagger"
echo "Database:  localhost:5432"
echo "=========================================="
echo ""

# Step 9: Display logs command
echo "📋 View Logs"
echo "=========================================="
echo "All services:"
echo "  docker compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f"
echo ""
echo "Backend only:"
echo "  docker compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f backend"
echo ""
echo "Frontend only:"
echo "  docker compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f frontend"
echo ""
echo "Database only:"
echo "  docker compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f postgres"
echo "=========================================="
echo ""

# Step 10: Health check results
echo "📋 Health Check Results"
echo "=========================================="

BACKEND_HEALTH=$(curl -s http://localhost:8080/health 2>/dev/null || echo "FAILED")
FRONTEND_HEALTH=$(curl -s http://localhost:3000/health 2>/dev/null || echo "FAILED")

echo "Backend:  $BACKEND_HEALTH"
echo "Frontend: $FRONTEND_HEALTH"
echo "=========================================="
echo ""

echo "✅ Deployment completed successfully!"
echo ""
echo "Next steps:"
echo "1. Verify all services are running: docker compose ps"
echo "2. Check logs: docker compose logs -f"
echo "3. Test API: curl http://localhost:8080/health"
echo "4. Open browser: http://localhost:3000"
echo ""
