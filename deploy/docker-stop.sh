#!/bin/bash

# ============================================
# ITSM Platform - Docker Stop Script
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

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Stop services
stop_services() {
    print_header "Stopping Services"
    
    docker-compose down
    
    if [ $? -eq 0 ]; then
        print_success "Services stopped successfully"
    else
        print_warning "Some services may not have stopped cleanly"
    fi
}

# Optional: Remove volumes
remove_volumes() {
    print_header "Remove Volumes?"
    
    read -p "Do you want to remove all data volumes? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose down -v
        print_success "Volumes removed"
    else
        print_warning "Volumes preserved"
    fi
}

# Show remaining containers
show_status() {
    print_header "Remaining Containers"
    
    docker ps -a | grep itsm || echo "No ITSM containers found"
}

# Main execution
main() {
    print_header "ITSM Platform - Docker Stop"
    
    stop_services
    remove_volumes
    show_status
    
    print_header "Stop Complete!"
}

# Run main function
main
