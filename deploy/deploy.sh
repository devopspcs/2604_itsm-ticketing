#!/bin/bash
# PCS Ticketing System - Production Deployment Script
set -e

echo "=== PCS Ticketing System - Production Deploy ==="

# 1. Pull latest code
echo "Pulling latest code..."
git pull origin main

# 2. Build and start containers
echo "Building and starting containers..."
docker compose -f docker-compose.prod.yml --env-file .env.prod build --no-cache
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 3. Wait for backend to be ready
echo "Waiting for backend..."
sleep 10

# 4. Health check
echo "Checking health..."
curl -s http://127.0.0.1:8080/health |45.77.173.144 python3 -m json.tool

echo ""
echo "=== Deployment Complete ==="
echo "Frontend: http://127.0.0.1:3000"
echo "Backend:  http://127.0.0.1:8080"
echo ""
echo "Make sure Apache is configured and running:"
echo "  sudo a2enmod proxy proxy_http rewrite"
echo "  sudo cp deploy/apache/itsm.conf /etc/apache2/sites-available/"
echo "  sudo a2ensite itsm.conf"
echo "  sudo systemctl reload apache2"
