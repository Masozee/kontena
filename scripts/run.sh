#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Starting Kontena CRM API ===${NC}"

# Check if the database exists
echo -e "${BLUE}Checking if database exists...${NC}"
DB_EXISTS=$(psql -U postgres -lqt | cut -d \| -f 1 | grep -w kontena | wc -l)

if [ "$DB_EXISTS" -eq "0" ]; then
    echo -e "${BLUE}Database does not exist. Creating...${NC}"
    psql -U postgres -c "CREATE DATABASE kontena;"
    echo -e "${GREEN}Database created successfully.${NC}"
else
    echo -e "${GREEN}Database already exists.${NC}"
fi

# Run the application
echo -e "${BLUE}Starting the application...${NC}"
go run cmd/api/main.go 