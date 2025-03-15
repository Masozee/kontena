package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/handlers"
	"github.com/Masozee/kontena/api/middleware"
)

// @title Project Management API
// @version 1.0
// @description Multi-tenant project management API built with Golang and Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Initialize database
	database.InitDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Project Management API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := app.Group("/api/v1")

	// Public routes (no tenant middleware)
	tenants := api.Group("/tenants")
	tenants.Get("/", handlers.GetTenants)
	tenants.Get("/:id", handlers.GetTenant)
	tenants.Post("/", handlers.CreateTenant)
	tenants.Put("/:id", handlers.UpdateTenant)
	tenants.Delete("/:id", handlers.DeleteTenant)

	// Protected routes (with tenant middleware)
	api.Use(middleware.TenantMiddleware())

	// Project routes
	projects := api.Group("/projects")
	projects.Get("/", handlers.GetProjects)
	projects.Get("/:id", handlers.GetProject)
	projects.Get("/:id/details", handlers.GetProjectWithDetails)
	projects.Post("/", handlers.CreateProject)
	projects.Patch("/:id", handlers.UpdateProject)
	projects.Delete("/:id", handlers.DeleteProject)

	// Person routes
	people := api.Group("/people")
	people.Get("/", handlers.GetPeople)
	people.Get("/:id", handlers.GetPerson)
	people.Post("/", handlers.CreatePerson)
	people.Patch("/:id", handlers.UpdatePerson)
	people.Delete("/:id", handlers.DeletePerson)

	// KPI routes
	kpis := api.Group("/kpis")
	kpis.Get("/:id", handlers.GetKPI)
	kpis.Patch("/:id", handlers.UpdateKPI)
	kpis.Delete("/:id", handlers.DeleteKPI)

	// Project KPI routes
	projectKpis := api.Group("/projects/:project_id/kpis")
	projectKpis.Get("/", handlers.GetKPIs)
	projectKpis.Post("/", handlers.CreateKPI)

	// Task routes
	tasks := api.Group("/tasks")
	tasks.Get("/", handlers.GetTasks)
	tasks.Get("/:id", handlers.GetTask)
	tasks.Post("/", handlers.CreateTask)
	tasks.Patch("/:id", handlers.UpdateTask)
	tasks.Delete("/:id", handlers.DeleteTask)

	// Project Task routes
	projectTasks := api.Group("/projects/:project_id/tasks")
	projectTasks.Get("/", handlers.GetTasks)
	projectTasks.Post("/", handlers.CreateTask)

	// Asset Management Routes

	// Asset Category routes
	assetCategories := api.Group("/asset-categories")
	assetCategories.Get("/", handlers.GetAssetCategories)
	assetCategories.Get("/:id", handlers.GetAssetCategory)
	assetCategories.Post("/", handlers.CreateAssetCategory)
	assetCategories.Put("/:id", handlers.UpdateAssetCategory)
	assetCategories.Delete("/:id", handlers.DeleteAssetCategory)

	// Asset routes
	assets := api.Group("/assets")
	assets.Get("/", handlers.GetAssets)
	assets.Get("/:id", handlers.GetAsset)
	assets.Post("/", handlers.CreateAsset)
	assets.Put("/:id", handlers.UpdateAsset)
	assets.Delete("/:id", handlers.DeleteAsset)

	// Location routes
	locations := api.Group("/locations")
	locations.Get("/", handlers.GetLocations)
	locations.Get("/:id", handlers.GetLocation)
	locations.Post("/", handlers.CreateLocation)
	locations.Put("/:id", handlers.UpdateLocation)
	locations.Delete("/:id", handlers.DeleteLocation)

	// Vendor routes
	vendors := api.Group("/vendors")
	vendors.Get("/", handlers.GetVendors)
	vendors.Get("/:id", handlers.GetVendor)
	vendors.Post("/", handlers.CreateVendor)
	vendors.Put("/:id", handlers.UpdateVendor)
	vendors.Delete("/:id", handlers.DeleteVendor)

	// Procurement Request routes
	procurementRequests := api.Group("/procurement-requests")
	procurementRequests.Get("/", handlers.GetProcurementRequests)
	procurementRequests.Get("/:id", handlers.GetProcurementRequest)
	procurementRequests.Post("/", handlers.CreateProcurementRequest)
	procurementRequests.Put("/:id", handlers.UpdateProcurementRequest)
	procurementRequests.Delete("/:id", handlers.DeleteProcurementRequest)

	// Asset Assignment routes
	assetAssignments := api.Group("/asset-assignments")
	assetAssignments.Get("/", handlers.GetAssetAssignments)
	assetAssignments.Get("/:id", handlers.GetAssetAssignment)
	assetAssignments.Post("/", handlers.CreateAssetAssignment)
	assetAssignments.Put("/:id", handlers.UpdateAssetAssignment)
	assetAssignments.Delete("/:id", handlers.DeleteAssetAssignment)

	// Maintenance Record routes
	maintenanceRecords := api.Group("/maintenance-records")
	maintenanceRecords.Get("/", handlers.GetMaintenanceRecords)
	maintenanceRecords.Get("/:id", handlers.GetMaintenanceRecord)
	maintenanceRecords.Post("/", handlers.CreateMaintenanceRecord)
	maintenanceRecords.Put("/:id", handlers.UpdateMaintenanceRecord)
	maintenanceRecords.Delete("/:id", handlers.DeleteMaintenanceRecord)

	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
