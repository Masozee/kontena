"Design a scalable, multi-tenant CRM API using Golang + Fiber, optimized for local development. The API should be discoverable without authentication initially and should support tenant-based data isolation. Key requirements include:

1️⃣ Multi-Tenancy Support:

Use shared-database, row-based isolation (tenant_id field on all tables).
API requests should support tenant filtering via headers or query parameters.
Tenant metadata should include company name, plan, and status.
2️⃣ User & Access Management (No Authentication Yet):

Users belong to a tenant and have roles (admin, sales, support).
Implement category-based access control (ABAC) to restrict access to leads.
Extendable user profiles with JSON fields for custom attributes.
3️⃣ Discoverable API (Swagger/OpenAPI):

Implement fiber-swagger for API documentation.
Ensure each endpoint is documented with request & response formats.
4️⃣ Database Choice & Schema:

Use PostgreSQL with GORM as ORM.
Support JSONB fields for dynamic user profiles.
Use migrations for database schema management.
5️⃣ API & Middleware:

Implement multi-tenant middleware to extract tenant_id dynamically.
No authentication required initially, but structure it to allow easy addition later.
6️⃣ Local Development Setup:

Run locally using Docker + Docker Compose (optional).
Include a Makefile or simple script for setting up the database.
Support .env-based configuration for database credentials.
7️⃣ Scalability Considerations (for future expansion):

Structure the codebase to allow authentication integration later (JWT/OAuth2).
Design API endpoints with pagination & filtering for handling large datasets.
Optimize queries to ensure tenant-based row filtering is efficient.
🚀 Generate a full system architecture, including recommended technologies, database schema, and API design for handling user roles, profiles, and category-based access control efficiently. Provide a Golang + Fiber implementation with Swagger integration for discoverability."

Expected Output:
✅ High-Level Architecture Diagram (Fiber API, DB, Multi-Tenant Structure)
✅ Database Schema (Tables for Tenant, User, Category, Leads, User Access)
✅ API Endpoints (/tenants, /users, /categories, /leads)
✅ Swagger Integration (OpenAPI Docs with fiber-swagger)
✅ Middleware for Multi-Tenancy (Header/Query Param Based)
✅ Local Setup Guide (Docker + .env Config)