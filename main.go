package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/your-username/project-management/database"
	"github.com/your-username/project-management/handlers"
	"github.com/your-username/project-management/middleware"
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
	kpis.Get("/", handlers.GetKPIs)
	kpis.Get("/:id", handlers.GetKPI)
	kpis.Post("/", handlers.CreateKPI)
	kpis.Patch("/:id", handlers.UpdateKPI)
	kpis.Delete("/:id", handlers.DeleteKPI)

	// Task routes
	tasks := api.Group("/tasks")
	tasks.Get("/", handlers.GetTasks)
	tasks.Get("/:id", handlers.GetTask)
	tasks.Post("/", handlers.CreateTask)
	tasks.Patch("/:id", handlers.UpdateTask)
	tasks.Delete("/:id", handlers.DeleteTask)

	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
