package handlers_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Masozee/kontena/api/internal/database"
	"github.com/Masozee/kontena/api/internal/handlers"
	"github.com/Masozee/kontena/api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB sets up an in-memory SQLite database for testing
func setupTestDB() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database")
	}

	// Auto migrate the models
	database.DB.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Category{},
		&models.Lead{},
	)

	// Create a test tenant
	tenant := models.Tenant{
		Name:   "Test Tenant",
		Plan:   "Basic",
		Status: "Active",
	}
	database.DB.Create(&tenant)
}

// setupApp sets up a Fiber app for testing
func setupApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestGetTenants(t *testing.T) {
	// Setup
	setupTestDB()
	app := setupApp()
	app.Get("/tenants", handlers.GetTenants)

	// Test
	req := httptest.NewRequest("GET", "/tenants", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var tenants []models.Tenant
	err = json.Unmarshal(body, &tenants)
	assert.NoError(t, err)
	assert.Len(t, tenants, 1)
	assert.Equal(t, "Test Tenant", tenants[0].Name)
}

func TestCreateTenant(t *testing.T) {
	// Setup
	setupTestDB()
	app := setupApp()
	app.Post("/tenants", handlers.CreateTenant)

	// Test data
	tenantData := `{"name":"New Tenant","plan":"Enterprise","status":"Active"}`

	// Test
	req := httptest.NewRequest("POST", "/tenants", strings.NewReader(tenantData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var tenant models.Tenant
	err = json.Unmarshal(body, &tenant)
	assert.NoError(t, err)
	assert.Equal(t, "New Tenant", tenant.Name)
	assert.Equal(t, "Enterprise", tenant.Plan)
	assert.Equal(t, "Active", tenant.Status)

	// Verify tenant was created in the database
	var count int64
	database.DB.Model(&models.Tenant{}).Count(&count)
	assert.Equal(t, int64(2), count) // 1 from setup + 1 from test
}

func TestGetTenant(t *testing.T) {
	// Setup
	setupTestDB()
	app := setupApp()
	app.Get("/tenants/:id", handlers.GetTenant)

	// Test
	req := httptest.NewRequest("GET", "/tenants/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var tenant models.Tenant
	err = json.Unmarshal(body, &tenant)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), tenant.ID)
	assert.Equal(t, "Test Tenant", tenant.Name)
}

func TestUpdateTenant(t *testing.T) {
	// Setup
	setupTestDB()
	app := setupApp()
	app.Put("/tenants/:id", handlers.UpdateTenant)

	// Test data
	tenantData := `{"name":"Updated Tenant","plan":"Premium","status":"Active"}`

	// Test
	req := httptest.NewRequest("PUT", "/tenants/1", strings.NewReader(tenantData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var tenant models.Tenant
	err = json.Unmarshal(body, &tenant)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), tenant.ID)
	assert.Equal(t, "Updated Tenant", tenant.Name)
	assert.Equal(t, "Premium", tenant.Plan)

	// Verify tenant was updated in the database
	var dbTenant models.Tenant
	database.DB.First(&dbTenant, 1)
	assert.Equal(t, "Updated Tenant", dbTenant.Name)
	assert.Equal(t, "Premium", dbTenant.Plan)
}

func TestDeleteTenant(t *testing.T) {
	// Setup
	setupTestDB()
	app := setupApp()
	app.Delete("/tenants/:id", handlers.DeleteTenant)

	// Test
	req := httptest.NewRequest("DELETE", "/tenants/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	// Verify tenant was deleted in the database
	var count int64
	database.DB.Model(&models.Tenant{}).Count(&count)
	assert.Equal(t, int64(0), count)
}
