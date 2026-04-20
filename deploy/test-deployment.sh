#!/bin/bash

# ============================================
# ITSM Platform - Deployment Test Script
# ============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Test backend health
test_backend() {
    print_header "Testing Backend"
    
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
    
    if [ "$response" = "200" ]; then
        print_success "Backend health check passed (HTTP $response)"
    else
        print_error "Backend health check failed (HTTP $response)"
        return 1
    fi
}

# Test frontend health
test_frontend() {
    print_header "Testing Frontend"
    
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/health)
    
    if [ "$response" = "200" ]; then
        print_success "Frontend health check passed (HTTP $response)"
    else
        print_error "Frontend health check failed (HTTP $response)"
        return 1
    fi
}

# Test database connection
test_database() {
    print_header "Testing Database"
    
    if docker-compose exec -T postgres pg_isready -U itsm &> /dev/null; then
        print_success "Database connection successful"
    else
        print_error "Database connection failed"
        return 1
    fi
}

# Test API endpoints
test_api() {
    print_header "Testing API Endpoints"
    
    # Test login endpoint
    response=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/v1/auth/login \
        -H "Content-Type: application/json" \
        -d '{"email":"test@example.com","password":"test"}')
    
    if [ "$response" = "401" ] || [ "$response" = "400" ]; then
        print_success "Login endpoint accessible (HTTP $response)"
    else
        print_warning "Login endpoint returned unexpected status (HTTP $response)"
    fi
    
    # Test health endpoint
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
    
    if [ "$response" = "200" ]; then
        print_success "Health endpoint accessible (HTTP $response)"
    else
        print_error "Health endpoint failed (HTTP $response)"
        return 1
    fi
}

# Test container status
test_containers() {
    print_header "Testing Container Status"
    
    # Check if all containers are running
    containers=$(docker-compose ps --services --filter "status=running")
    
    for service in postgres backend frontend; do
        if echo "$containers" | grep -q "$service"; then
            print_success "Container '$service' is running"
        else
            print_error "Container '$service' is not running"
            return 1
        fi
    done
}

# Test network connectivity
test_network() {
    print_header "Testing Network Connectivity"
    
    # Test backend can reach database
    if docker-compose exec -T backend wget --quiet --tries=1 --spider http://localhost:8080/health &> /dev/null; then
        print_success "Backend is accessible from network"
    else
        print_error "Backend is not accessible from network"
        return 1
    fi
}

# Show performance metrics
show_metrics() {
    print_header "Performance Metrics"
    
    echo "Container Resource Usage:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
}

# Show service information
show_info() {
    print_header "Service Information"
    
    echo "Running Services:"
    docker-compose ps
    
    echo ""
    echo "Network:"
    docker network ls | grep itsm
    
    echo ""
    echo "Volumes:"
    docker volume ls | grep itsm
}

# Main execution
main() {
    print_header "ITSM Platform - Deployment Test"
    
    # Check if services are running
    if ! docker-compose ps | grep -q "Up"; then
        print_error "Services are not running. Please start them first with: docker-compose up -d"
        exit 1
    fi
    
    # Run tests
    test_containers || exit 1
    test_database || exit 1
    test_backend || exit 1
    test_frontend || exit 1
    test_api || exit 1
    test_network || exit 1
    
    # Show information
    show_info
    show_metrics
    
    print_header "All Tests Passed!"
    print_success "ITSM Platform deployment is working correctly"
}

# Run main function
main
