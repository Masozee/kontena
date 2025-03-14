
"Create a comprehensive **self-testing strategy** and **dummy data generation script** for a multi-tenant CRM API built with **Golang + Fiber + PostgreSQL**. The system should support multi-tenancy via `tenant_id`, and testing should include inserting, retrieving, and validating data integrity. The tests should ensure API endpoints work as expected before authentication is added. 

---

### **1️⃣ Dummy Data Generation**  
**✅ Generate SQL/Seed Scripts to Populate Database with:**  
- **Tenants:**  
  - `tenant_id`, `name`, `plan`, `status`  
  - Example: `{ "tenant_id": 1, "name": "Acme Corp", "plan": "Enterprise", "status": "Active" }`  

- **Users (Belonging to Tenants):**  
  - `user_id`, `tenant_id`, `name`, `role` (`admin`, `sales`, `support`), `email`, `profile (JSONB)`  
  - Example: `{ "user_id": 1, "tenant_id": 1, "name": "Alice", "role": "admin", "email": "alice@acme.com" }`  

- **Categories (Access Control):**  
  - `category_id`, `tenant_id`, `name`, `permissions`  
  - Example: `{ "category_id": 1, "tenant_id": 1, "name": "VIP Clients", "permissions": ["read", "write"] }`  

- **Leads (CRM Data):**  
  - `lead_id`, `tenant_id`, `name`, `contact`, `category_id`, `assigned_to`  
  - Example: `{ "lead_id": 1, "tenant_id": 1, "name": "John Doe", "contact": "123456789", "category_id": 1, "assigned_to": 2 }`  

---

### **2️⃣ API Self-Testing (Using `Go Test` & `Postman` Requests)**  
**✅ Define API Tests to Verify:**  
- **Tenant-Specific Filtering:**  
  - Ensure `GET /users?tenant_id=1` returns only **users for that tenant**.  
  - Ensure `GET /leads?tenant_id=1` does not return **leads from tenant 2**.  

- **CRUD Operations for Key Endpoints:**  
  - **Create, Read, Update, Delete** tests for `/tenants`, `/users`, `/categories`, `/leads`.  
  - Ensure correct **HTTP response codes** (`201`, `200`, `400`, `404`).  

- **Category-Based Access Control (ABAC):**  
  - Ensure **users with `sales` role** can access leads.  
  - Ensure **users with `support` role** cannot update `VIP Clients`.  

- **Multi-Tenancy Edge Cases:**  
  - Test API responses **without `tenant_id`** (should return error).  
  - Test cross-tenant access **(tenant A cannot access tenant B’s data)**.  

---

### **3️⃣ Script for Populating Database (Go + SQL)**  
**✅ Generate a Go script to:**  
- Connect to **PostgreSQL locally**.  
- Insert **dummy tenants, users, categories, and leads**.  
- Log inserted records for validation.  

---

### **🚀 Expected Output**  
- ✅ **Pre-filled Database with Dummy CRM Data**  
- ✅ **Automated API Tests for Multi-Tenant Filtering**  
- ✅ **Sample Requests for Postman or Curl**  
- ✅ **Golang Script for Dummy Data Insertion**  

---

