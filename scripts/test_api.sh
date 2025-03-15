#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# API base URL
API_URL="http://localhost:3000/api/v1"

# Function to make API requests and display results
function make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local tenant_id=$4
    
    echo -e "${YELLOW}${BOLD}$method $endpoint${NC}"
    
    if [ "$method" == "GET" ]; then
        if [ ! -z "$tenant_id" ]; then
            curl -s -X $method -H "X-Tenant-ID: $tenant_id" $API_URL$endpoint | jq '.' 2>/dev/null || echo "Error parsing JSON"
        else
            curl -s -X $method $API_URL$endpoint | jq '.' 2>/dev/null || echo "Error parsing JSON"
        fi
    else
        if [ ! -z "$tenant_id" ]; then
            curl -s -X $method -H "Content-Type: application/json" -H "X-Tenant-ID: $tenant_id" -d "$data" $API_URL$endpoint | jq '.' 2>/dev/null || echo "Error parsing JSON"
        else
            curl -s -X $method -H "Content-Type: application/json" -d "$data" $API_URL$endpoint | jq '.' 2>/dev/null || echo "Error parsing JSON"
        fi
    fi
    
    echo ""
}

echo -e "${GREEN}${BOLD}Testing Multi-Tenant Project Management API${NC}"
echo "========================================"
echo ""

# Test tenant endpoints
echo -e "${GREEN}${BOLD}Testing Tenant Endpoints${NC}"
echo "------------------------"
make_request "GET" "/tenants"
make_request "GET" "/tenants/1"
make_request "POST" "/tenants" '{"name":"Test Tenant","description":"A test tenant","plan":"free","status":"active","domain":"test.example.com"}'

# Test project endpoints for tenant 1
echo -e "${GREEN}${BOLD}Testing Project Endpoints${NC}"
echo "------------------------"
make_request "GET" "/projects" "" "1"
make_request "GET" "/projects/1" "" "1"
make_request "GET" "/projects/1/details" "" "1"
make_request "POST" "/projects" '{"name":"New Test Project","description":"A test project","budget":25000,"status":"planning","start_date":"2023-01-01T00:00:00Z"}' "1"

# Test person endpoints for tenant 1
echo -e "${GREEN}${BOLD}Testing Person Endpoints${NC}"
echo "------------------------"
make_request "GET" "/people" "" "1"
make_request "GET" "/people/1" "" "1"
make_request "POST" "/people" '{"name":"New Test Person","email":"test.person@example.com","role":"Developer","position":"Junior Developer","phone":"555-1234"}' "1"

# Test task endpoints for tenant 1, project 1
echo -e "${GREEN}${BOLD}Testing Task Endpoints${NC}"
echo "------------------------"
make_request "GET" "/projects/1/tasks" "" "1"
make_request "GET" "/tasks/1" "" "1"
make_request "POST" "/projects/1/tasks" '{"title":"New Test Task","description":"A test task","status":"todo"}' "1"

# Test KPI endpoints for tenant 1, project 1
echo -e "${GREEN}${BOLD}Testing KPI Endpoints${NC}"
echo "------------------------"
make_request "GET" "/projects/1/kpis" "" "1"
make_request "GET" "/kpis/1" "" "1"
make_request "POST" "/projects/1/kpis" '{"description":"New Test KPI","target_value":100,"current_value":50,"unit":"%"}' "1"

# Test milestone endpoints for tenant 1, project 1
echo -e "${GREEN}${BOLD}Testing Milestone Endpoints${NC}"
echo "------------------------"
make_request "GET" "/projects/1/milestones" "" "1"
make_request "GET" "/milestones/1" "" "1"
make_request "POST" "/projects/1/milestones" '{"title":"New Test Milestone","description":"A test milestone","due_date":"2023-12-31T00:00:00Z","status":"planned"}' "1"

# Test risk endpoints for tenant 1, project 1
echo -e "${GREEN}${BOLD}Testing Risk Endpoints${NC}"
echo "------------------------"
make_request "GET" "/projects/1/risks" "" "1"
make_request "GET" "/risks/1" "" "1"
make_request "POST" "/projects/1/risks" '{"description":"New Test Risk","impact":"Medium","probability":"Low","mitigation":"Test mitigation strategy","status":"identified"}' "1"

# Test time tracking endpoints for tenant 1
echo -e "${GREEN}${BOLD}Testing Time Tracking Endpoints${NC}"
echo "------------------------"
make_request "GET" "/time-entries" "" "1"
make_request "GET" "/time-entries/1" "" "1"
make_request "POST" "/tasks/1/time-entries" '{"person_id":1,"hours":2.5,"date":"2023-01-01T00:00:00Z"}' "1"

# Test multi-tenancy isolation
echo -e "${GREEN}${BOLD}Testing Multi-Tenancy Isolation${NC}"
echo "------------------------"
echo -e "${YELLOW}${BOLD}Attempting to access tenant 2's projects with tenant 1's ID${NC}"
make_request "GET" "/projects" "" "1"
echo -e "${YELLOW}${BOLD}Attempting to access tenant 1's projects with tenant 2's ID${NC}"
make_request "GET" "/projects" "" "2"

echo -e "${GREEN}${BOLD}API Testing Complete${NC}" 