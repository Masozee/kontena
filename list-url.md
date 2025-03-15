# API Endpoints Reference

This document lists all the working API endpoints for the multi-tenant project management system. All examples use Tenant ID 1.

## Authentication

All protected endpoints require a tenant ID to be provided in one of the following ways:
- Header: `X-Tenant-ID: 1`
- Header: `tenant_id: 1`
- Query parameter: `?tenant_id=1`

## Tenant Endpoints (Public)

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/tenants | Get all tenants |
| GET | http://localhost:3000/api/v1/tenants/1 | Get tenant by ID |
| POST | http://localhost:3000/api/v1/tenants | Create a new tenant |
| PUT | http://localhost:3000/api/v1/tenants/1 | Update a tenant |
| DELETE | http://localhost:3000/api/v1/tenants/1 | Delete a tenant |

## Project Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/projects | Get all projects for tenant |
| GET | http://localhost:3000/api/v1/projects/16 | Get project by ID |
| GET | http://localhost:3000/api/v1/projects/16/details | Get project with details |
| POST | http://localhost:3000/api/v1/projects | Create a new project |
| PATCH | http://localhost:3000/api/v1/projects/16 | Update a project |
| DELETE | http://localhost:3000/api/v1/projects/16 | Delete a project |

## Person Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/people | Get all people for tenant |
| GET | http://localhost:3000/api/v1/people/16 | Get person by ID |
| POST | http://localhost:3000/api/v1/people | Create a new person |
| PATCH | http://localhost:3000/api/v1/people/16 | Update a person |
| DELETE | http://localhost:3000/api/v1/people/16 | Delete a person |

## Task Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/tasks | Get all tasks for tenant |
| GET | http://localhost:3000/api/v1/tasks/47 | Get task by ID |
| POST | http://localhost:3000/api/v1/tasks | Create a new task |
| PATCH | http://localhost:3000/api/v1/tasks/47 | Update a task |
| DELETE | http://localhost:3000/api/v1/tasks/47 | Delete a task |

## Project-specific Task Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/projects/16/tasks | Get all tasks for a project |
| POST | http://localhost:3000/api/v1/projects/16/tasks | Create a new task for a project |

## KPI Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/kpis/1 | Get KPI by ID |
| PATCH | http://localhost:3000/api/v1/kpis/1 | Update a KPI |
| DELETE | http://localhost:3000/api/v1/kpis/1 | Delete a KPI |

## Project-specific KPI Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | http://localhost:3000/api/v1/projects/16/kpis | Get all KPIs for a project |
| POST | http://localhost:3000/api/v1/projects/16/kpis | Create a new KPI for a project |

## Example cURL Commands

### Get all projects for tenant 1
```bash
curl -s -H "X-Tenant-ID: 1" http://localhost:3000/api/v1/projects | jq
```

### Create a new project for tenant 1
```bash
curl -s -X POST -H "X-Tenant-ID: 1" -H "Content-Type: application/json" \
  -d '{"name":"Test Project", "description":"This is a test project", "status":"active"}' \
  http://localhost:3000/api/v1/projects | jq
```

### Create a new task for project 16
```bash
curl -s -X POST -H "X-Tenant-ID: 1" -H "Content-Type: application/json" \
  -d '{"title":"Test Task", "description":"This is a test task", "status":"todo"}' \
  http://localhost:3000/api/v1/projects/16/tasks | jq
```

### Update a task
```bash
curl -s -X PATCH -H "X-Tenant-ID: 1" -H "Content-Type: application/json" \
  -d '{"title":"Updated Task", "description":"This is an updated task", "status":"in_progress"}' \
  http://localhost:3000/api/v1/tasks/47 | jq
```

### Delete a task
```bash
curl -s -X DELETE -H "X-Tenant-ID: 1" http://localhost:3000/api/v1/tasks/47 | jq
``` 