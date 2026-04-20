#!/bin/bash

# Comprehensive Jira Setup Testing Script
# This script tests all components and fixes issues

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Jira Project Board - Comprehensive Test${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-itsm}
DB_PASSWORD=${DB_PASSWORD:-itsm}
DB_NAME=${DB_NAME:-itsm}
API_URL=${API_URL:-http://localhost:8080/api/v1}
FRONTEND_URL=${FRONTEND_URL:-http://localhost:3000}

print_status() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}✓ $2${NC}"
  else
    echo -e "${RED}✗ $2${NC}"
  fi
}

print_info() {
  echo -e "${BLUE}ℹ $1${NC}"
}

print_warning() {
  echo -e "${YELLOW}⚠ $1${NC}"
}

# Test 1: Check if Docker containers are running
echo -e "${BLUE}1. Checking Docker Containers...${NC}"
if docker ps | grep -q "itsm-postgres"; then
  print_status 0 "PostgreSQL container running"
else
  print_warning "PostgreSQL container not running"
  echo "Starting Docker containers..."
  docker-compose up -d
  sleep 5
fi

if docker ps | grep -q "itsm-backend"; then
  print_status 0 "Backend container running"
else
  print_warning "Backend container not running"
fi

if docker ps | grep -q "itsm-frontend"; then
  print_status 0 "Frontend container running"
else
  print_warning "Frontend container not running"
fi

# Test 2: Check database connection
echo -e "\n${BLUE}2. Testing Database Connection...${NC}"
if docker exec itsm-postgres pg_isready -U $DB_USER > /dev/null 2>&1; then
  print_status 0 "Database connection successful"
else
  print_status 1 "Database connection failed"
  exit 1
fi

# Test 3: Check if Jira tables exist
echo -e "\n${BLUE}3. Checking Database Schema...${NC}"

check_table() {
  if docker exec itsm-postgres psql -U $DB_USER -d $DB_NAME -c "SELECT 1 FROM information_schema.tables WHERE table_name='$1'" 2>/dev/null | grep -q "1"; then
    print_status 0 "$1 table exists"
    return 0
  else
    print_warning "$1 table not found"
    return 1
  fi
}

check_table "issue_types"
check_table "workflows"
check_table "sprints"
check_table "custom_fields"
check_table "labels"

# Test 4: Check test data
echo -e "\n${BLUE}4. Checking Test Data...${NC}"

# Check issue types
count=$(docker exec itsm-postgres psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM issue_types" 2>/dev/null | tr -d ' ')
if [ "$count" -gt 0 ]; then
  print_status 0 "Issue types exist ($count rows)"
else
  print_warning "No issue types found - need to insert test data"
fi

# Check projects
count=$(docker exec itsm-postgres psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM projects WHERE name LIKE '%Jira%'" 2>/dev/null | tr -d ' ')
if [ "$count" -gt 0 ]; then
  print_status 0 "Test project exists"
else
  print_warning "Test project not found"
fi

# Check sprints
count=$(docker exec itsm-postgres psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM sprints WHERE status='Active'" 2>/dev/null | tr -d ' ')
if [ "$count" -gt 0 ]; then
  print_status 0 "Active sprint exists"
else
  print_warning "No active sprint found"
fi

# Test 5: Test API endpoints
echo -e "\n${BLUE}5. Testing API Endpoints...${NC}"

# Get auth token
print_info "Getting auth token..."
TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  print_warning "Could not get auth token - trying without auth"
  TOKEN=""
else
  print_status 0 "Auth token obtained"
fi

# Get projects
print_info "Testing GET /projects..."
if curl -s -f "$API_URL/projects" \
  -H "Authorization: Bearer $TOKEN" > /dev/null 2>&1; then
  print_status 0 "GET /projects works"
else
  print_warning "GET /projects failed"
fi

# Get issue types (need project ID)
print_info "Getting project ID..."
PROJECT_ID=$(curl -s "$API_URL/projects" \
  -H "Authorization: Bearer $TOKEN" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -n "$PROJECT_ID" ]; then
  print_status 0 "Project ID: $PROJECT_ID"
  
  # Test issue types endpoint
  print_info "Testing GET /projects/{id}/issue-types..."
  if curl -s -f "$API_URL/projects/$PROJECT_ID/issue-types" \
    -H "Authorization: Bearer $TOKEN" > /dev/null 2>&1; then
    print_status 0 "GET /projects/{id}/issue-types works"
  else
    print_warning "GET /projects/{id}/issue-types failed"
  fi
  
  # Test workflow endpoint
  print_info "Testing GET /projects/{id}/workflow..."
  if curl -s -f "$API_URL/projects/$PROJECT_ID/workflow" \
    -H "Authorization: Bearer $TOKEN" > /dev/null 2>&1; then
    print_status 0 "GET /projects/{id}/workflow works"
  else
    print_warning "GET /projects/{id}/workflow failed"
  fi
  
  # Test active sprint endpoint
  print_info "Testing GET /projects/{id}/sprints/active..."
  SPRINT_RESPONSE=$(curl -s "$API_URL/projects/$PROJECT_ID/sprints/active" \
    -H "Authorization: Bearer $TOKEN")
  if echo "$SPRINT_RESPONSE" | grep -q "id"; then
    print_status 0 "GET /projects/{id}/sprints/active works"
  else
    print_warning "GET /projects/{id}/sprints/active returned empty or error"
    echo "Response: $SPRINT_RESPONSE"
  fi
else
  print_warning "Could not get project ID"
fi

# Test 6: Check frontend
echo -e "\n${BLUE}6. Testing Frontend...${NC}"

if curl -s -f "$FRONTEND_URL" > /dev/null 2>&1; then
  print_status 0 "Frontend is accessible"
else
  print_warning "Frontend is not accessible at $FRONTEND_URL"
fi

# Summary
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"

print_info "Database: $DB_HOST:$DB_PORT/$DB_NAME"
print_info "API URL: $API_URL"
print_info "Frontend URL: $FRONTEND_URL"

echo -e "\n${YELLOW}Recommended Next Steps:${NC}"
echo "1. Check backend logs: docker logs itsm-backend"
echo "2. Check frontend console: Browser DevTools (F12)"
echo "3. Check database: docker exec itsm-postgres psql -U $DB_USER -d $DB_NAME"
echo "4. If tables missing, run migrations manually"
echo "5. If test data missing, insert seed data"

echo -e "\n${GREEN}Test complete!${NC}\n"
