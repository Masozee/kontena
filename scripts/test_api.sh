#!/bin/bash

# Set the API base URL
API_URL="http://localhost:3000/api/v1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to make API requests and display results
function make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local tenant_id=$4
    
    echo -e "${BLUE}Testing: ${method} ${endpoint}${NC}"
    
    # Set tenant header if provided
    local headers=""
    if [ ! -z "$tenant_id" ]; then
        headers="-H \"X-Tenant-ID: ${tenant_id}\""
    fi
    
    # Make the request
    if [ "$method" == "GET" ]; then
        if [ ! -z "$tenant_id" ]; then
            response=$(curl -s -X ${method} -H "X-Tenant-ID: ${tenant_id}" ${API_URL}${endpoint})
        else
            response=$(curl -s -X ${method} ${API_URL}${endpoint})
        fi
    else
        if [ ! -z "$tenant_id" ]; then
            response=$(curl -s -X ${method} -H "Content-Type: application/json" -H "X-Tenant-ID: ${tenant_id}" -d "${data}" ${API_URL}${endpoint})
        else
            response=$(curl -s -X ${method} -H "Content-Type: application/json" -d "${data}" ${API_URL}${endpoint})
        fi
    fi
    
    # Check if response is valid JSON
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo -e "${GREEN}Response:${NC}"
        echo "$response" | jq .
        echo ""
    else
        echo -e "${RED}Error: Invalid JSON response${NC}"
        echo "$response"
        echo ""
    fi
}

echo -e "${BLUE}=== Testing Kontena CRM API ===${NC}"

# Test tenant endpoints
echo -e "${BLUE}=== Testing Tenant Endpoints ===${NC}"
make_request "GET" "/tenants"

# Create a new tenant
echo -e "${BLUE}Creating a new tenant...${NC}"
tenant_data='{
    "name": "Test Tenant",
    "plan": "Basic",
    "status": "Active"
}'
make_request "POST" "/tenants" "$tenant_data"

# Get the first tenant
make_request "GET" "/tenants/1"

# Test user endpoints with tenant filtering
echo -e "${BLUE}=== Testing User Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/users" "" "1"

# Create a new user for tenant 1
echo -e "${BLUE}Creating a new user for tenant 1...${NC}"
user_data='{
    "tenant_id": 1,
    "name": "Test User",
    "email": "test@example.com",
    "role": "admin",
    "profile": "{\"department\":\"IT\",\"location\":\"San Francisco\"}"
}'
make_request "POST" "/users" "$user_data" "1"

# Test category endpoints with tenant filtering
echo -e "${BLUE}=== Testing Category Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/categories" "" "1"

# Create a new category for tenant 1
echo -e "${BLUE}Creating a new category for tenant 1...${NC}"
category_data='{
    "tenant_id": 1,
    "name": "Test Category",
    "permissions": "[\"read\", \"write\"]"
}'
make_request "POST" "/categories" "$category_data" "1"

# Test lead endpoints with tenant filtering
echo -e "${BLUE}=== Testing Lead Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/leads" "" "1"

# Create a new lead for tenant 1
echo -e "${BLUE}Creating a new lead for tenant 1...${NC}"
lead_data='{
    "tenant_id": 1,
    "name": "Test Lead",
    "contact": "555-123-4567",
    "email": "lead@example.com",
    "category_id": 1,
    "assigned_to": 1,
    "status": "New",
    "notes": "Test lead created via API"
}'
make_request "POST" "/leads" "$lead_data" "1"

# Test staff endpoints with tenant filtering
echo -e "${BLUE}=== Testing Staff Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/staff" "" "1"

# Create a new staff member for tenant 1
echo -e "${BLUE}Creating a new staff member for tenant 1...${NC}"
staff_data='{
    "tenant_id": 1,
    "name": "Test Staff",
    "email": "staff@example.com",
    "role": "admin",
    "department": "IT",
    "position": "Manager",
    "profile": "{\"skills\":[\"Leadership\",\"Communication\"],\"education\":\"Bachelor\\\"s Degree\"}"
}'
make_request "POST" "/staff" "$staff_data" "1"

# Test archive endpoints with tenant filtering
echo -e "${BLUE}=== Testing Archive Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/archives" "" "1"

# Create a new archive for tenant 1
echo -e "${BLUE}Creating a new archive for tenant 1...${NC}"
archive_data='{
    "tenant_id": 1,
    "title": "Test Archive",
    "description": "Test archive document",
    "category_id": 1,
    "file_path": "/uploads/test.pdf",
    "file_type": "application/pdf",
    "file_size": 1024,
    "status": "active",
    "created_by_id": 1,
    "metadata": "{\"tags\":[\"important\",\"document\"],\"department\":\"Legal\"}"
}'
make_request "POST" "/archives" "$archive_data" "1"

# Test asset endpoints with tenant filtering
echo -e "${BLUE}=== Testing Asset Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/assets" "" "1"

# Create a new asset for tenant 1
echo -e "${BLUE}Creating a new asset for tenant 1...${NC}"
asset_data='{
    "tenant_id": 1,
    "name": "Test Asset",
    "description": "Test computer asset",
    "type": "computer",
    "serial_number": "SN-12345",
    "purchase_date": "2023-01-01T00:00:00Z",
    "purchase_price": 1000,
    "status": "available",
    "location": "Office 1",
    "assigned_to_id": 1,
    "notes": "Test asset created via API",
    "metadata": "{\"manufacturer\":\"Example Corp\",\"warranty\":\"2 years\"}"
}'
make_request "POST" "/assets" "$asset_data" "1"

# Test ticket endpoints with tenant filtering
echo -e "${BLUE}=== Testing Ticket Endpoints with Tenant Filtering ===${NC}"
make_request "GET" "/tickets" "" "1"

# Create a new ticket for tenant 1
echo -e "${BLUE}Creating a new ticket for tenant 1...${NC}"
ticket_data='{
    "tenant_id": 1,
    "title": "Test Ticket",
    "description": "Test support ticket",
    "priority": "medium",
    "status": "open",
    "created_by_id": 1,
    "assigned_to_id": 1,
    "due_date": "2023-12-31T00:00:00Z",
    "category_id": 1,
    "related_asset_id": 1,
    "notes": "Test ticket created via API",
    "metadata": "{\"source\":\"Email\",\"tags\":[\"hardware\",\"urgent\"]}"
}'
make_request "POST" "/tickets" "$ticket_data" "1"

# Test multi-tenancy by trying to access tenant 2's data with tenant 1's ID
echo -e "${BLUE}=== Testing Multi-Tenancy \(Cross-Tenant Access\) ===${NC}"
echo -e "${BLUE}Trying to access tenant 2's users with tenant 1's ID...${NC}"
make_request "GET" "/users" "" "1"

echo -e "${BLUE}=== API Testing Complete ===${NC}" 