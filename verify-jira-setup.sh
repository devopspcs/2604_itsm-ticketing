#!/bin/bash

# Jira-like Project Board Setup Verification Script
# This script verifies that all components are properly set up

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-itsm_db}
API_URL=${API_URL:-http://localhost:8080/api/v1}
FRONTEND_URL=${FRONTEND_URL:-http://localhost:5173}

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Jira-like Project Board Setup Verification${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Function to print status
print_status() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}✓ $2${NC}"
  else
    echo -e "${RED}✗ $2${NC}"
  fi
}

# Function to print info
print_info() {
  echo -e "${BLUE}ℹ $1${NC}"
}

# Function to print warning
print_warning() {
  echo -e "${YELLOW}⚠ $1${NC}"
}

# 1. Check Database Connection
echo -e "${BLUE}1. Checking Database Connection...${NC}"
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1" > /dev/null 2>&1; then
  print_status 0 "Database connection successful"
else
  print_status 1 "Database connection failed"
  exit 1
fi

# 2. Check Database Schema
echo -e "\n${BLUE}2. Checking Database Schema...${NC}"

# Check issue_types table
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) FROM issue_types" > /dev/null 2>&1; then
  count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM issue_types")
  print_status 0 "issue_types table exists ($count rows)"
else
  print_status 1 "issue_types table not found"
fi

# Check workflows table
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) FROM workflows" > /dev/null 2>&1; then
  count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM workflows")
  print_status 0 "workflows table exists ($count rows)"
else
  print_status 1 "workflows table not found"
fi

# Check sprints table
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) FROM sprints" > /dev/null 2>&1; then
  count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM sprints")
  print_status 0 "sprints table exists ($count rows)"
else
  print_status 1 "sprints table not found"
fi

# Check custom_fields table
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) FROM custom_fields" > /dev/null 2>&1; then
  count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM custom_fields")
  print_status 0 "custom_fields table exists ($count rows)"
else
  print_status 1 "custom_fields table not found"
fi

# Check labels table
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) FROM labels" > /dev/null 2>&1; then
  count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM labels")
  print_status 0 "labels table exists ($count rows)"
else
  print_status 1 "labels table not found"
fi

# 3. Check Test Data
echo -e "\n${BLUE}3. Checking Test Data...${NC}"

# Check predefined issue types
count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM issue_types WHERE name IN ('Bug', 'Task', 'Story', 'Epic', 'Sub-task')")
if [ "$count" -eq 5 ]; then
  print_status 0 "All 5 predefined issue types exist"
else
  print_warning "Only $count/5 predefined issue types found"
fi

# Check test project
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1 FROM projects WHERE name = 'Test Project - Jira Board'" > /dev/null 2>&1; then
  print_status 0 "Test project exists"
else
  print_warning "Test project not found"
fi

# Check active sprint
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1 FROM sprints WHERE status = 'Active'" > /dev/null 2>&1; then
  print_status 0 "Active sprint exists"
else
  print_warning "No active sprint found"
fi

# Check labels
count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM labels")
if [ "$count" -gt 0 ]; then
  print_status 0 "Labels exist ($count rows)"
else
  print_warning "No labels found"
fi

# Check custom fields
count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM custom_fields")
if [ "$count" -gt 0 ]; then
  print_status 0 "Custom fields exist ($count rows)"
else
  print_warning "No custom fields found"
fi

# 4. Check Backend API
echo -e "\n${BLUE}4. Checking Backend API...${NC}"

# Check if backend is running
if curl -s -f "$API_URL/health" > /dev/null 2>&1; then
  print_status 0 "Backend API is running"
else
  print_warning "Backend API health check failed (may not have /health endpoint)"
fi

# 5. Check Frontend
echo -e "\n${BLUE}5. Checking Frontend...${NC}"

# Check if frontend is running
if curl -s -f "$FRONTEND_URL" > /dev/null 2>&1; then
  print_status 0 "Frontend is running"
else
  print_warning "Frontend is not running at $FRONTEND_URL"
fi

# 6. Check File Structure
echo -e "\n${BLUE}6. Checking File Structure...${NC}"

# Check backend files
if [ -f "backend/migrations/000009_jira_features.up.sql" ]; then
  print_status 0 "Jira features migration exists"
else
  print_status 1 "Jira features migration not found"
fi

if [ -f "backend/migrations/000011_seed_jira_test_data.up.sql" ]; then
  print_status 0 "Test data migration exists"
else
  print_status 1 "Test data migration not found"
fi

# Check frontend files
if [ -f "frontend/src/services/jira.service.ts" ]; then
  print_status 0 "Jira service exists"
else
  print_status 1 "Jira service not found"
fi

if [ -f "frontend/src/pages/ProjectBoardPage.tsx" ]; then
  print_status 0 "ProjectBoardPage exists"
else
  print_status 1 "ProjectBoardPage not found"
fi

if [ -f "frontend/src/pages/SprintBoardPage.tsx" ]; then
  print_status 0 "SprintBoardPage exists"
else
  print_status 1 "SprintBoardPage not found"
fi

if [ -f "frontend/src/components/project/SprintBoard.tsx" ]; then
  print_status 0 "SprintBoard component exists"
else
  print_status 1 "SprintBoard component not found"
fi

# 7. Check Dependencies
echo -e "\n${BLUE}7. Checking Dependencies...${NC}"

# Check Go dependencies
if [ -f "backend/go.mod" ]; then
  if grep -q "github.com/lib/pq" backend/go.mod; then
    print_status 0 "PostgreSQL driver is in go.mod"
  else
    print_warning "PostgreSQL driver not found in go.mod"
  fi
else
  print_warning "go.mod not found"
fi

# Check npm dependencies
if [ -f "frontend/package.json" ]; then
  if grep -q "@dnd-kit/core" frontend/package.json; then
    print_status 0 "Drag-and-drop library is in package.json"
  else
    print_warning "Drag-and-drop library not found in package.json"
  fi
else
  print_warning "package.json not found"
fi

# 8. Summary
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}Verification Summary${NC}"
echo -e "${BLUE}========================================${NC}"

print_info "Database: $DB_HOST:$DB_PORT/$DB_NAME"
print_info "API URL: $API_URL"
print_info "Frontend URL: $FRONTEND_URL"

echo -e "\n${GREEN}Setup verification complete!${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo "1. Run database migrations: migrate -path backend/migrations -database \"postgresql://...\" up"
echo "2. Insert test data: psql -U $DB_USER -d $DB_NAME -f backend/migrations/000011_seed_jira_test_data.up.sql"
echo "3. Start backend: cd backend && go run cmd/main.go"
echo "4. Start frontend: cd frontend && npm run dev"
echo "5. Open browser: $FRONTEND_URL"
echo "6. Navigate to project board and verify functionality"

echo -e "\n${BLUE}For more details, see JIRA_BOARD_SETUP_GUIDE.md${NC}\n"
