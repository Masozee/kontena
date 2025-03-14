package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	_ "github.com/kontena/api/docs"
	"github.com/kontena/api/internal/database"
	"github.com/kontena/api/internal/handlers"
	"github.com/kontena/api/internal/middleware"
)

// @title Kontena CRM API
// @version 1.0
// @description Multi-tenant CRM API built with Golang and Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@kontena.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
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
		AppName: "Kontena CRM API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(middleware.TenantMiddleware())

	// Swagger documentation - updated configuration
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: true,
		Title:       "Kontena CRM API Documentation",
	}))

	// API routes
	api := app.Group("/api/v1")

	// Tenant routes
	tenants := api.Group("/tenants")
	tenants.Get("/", handlers.GetTenants)
	tenants.Get("/:id", handlers.GetTenant)
	tenants.Post("/", handlers.CreateTenant)
	tenants.Put("/:id", handlers.UpdateTenant)
	tenants.Delete("/:id", handlers.DeleteTenant)

	// User routes
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	// Category routes
	categories := api.Group("/categories")
	categories.Get("/", handlers.GetCategories)
	categories.Get("/:id", handlers.GetCategory)
	categories.Post("/", handlers.CreateCategory)
	categories.Put("/:id", handlers.UpdateCategory)
	categories.Delete("/:id", handlers.DeleteCategory)

	// Lead routes
	leads := api.Group("/leads")
	leads.Get("/", handlers.GetLeads)
	leads.Get("/:id", handlers.GetLead)
	leads.Post("/", handlers.CreateLead)
	leads.Put("/:id", handlers.UpdateLead)
	leads.Delete("/:id", handlers.DeleteLead)

	// Staff routes
	staff := api.Group("/staff")
	staff.Get("/", handlers.GetStaff)
	staff.Get("/:id", handlers.GetStaffMember)
	staff.Post("/", handlers.CreateStaffMember)
	staff.Patch("/:id", handlers.UpdateStaffMember)
	staff.Delete("/:id", handlers.DeleteStaffMember)

	// Archive routes
	archives := api.Group("/archives")
	archives.Get("/", handlers.GetArchives)
	archives.Get("/:id", handlers.GetArchive)
	archives.Post("/", handlers.CreateArchive)
	archives.Patch("/:id", handlers.UpdateArchive)
	archives.Delete("/:id", handlers.DeleteArchive)

	// Asset routes
	assets := api.Group("/assets")
	assets.Get("/", handlers.GetAssets)
	assets.Get("/:id", handlers.GetAsset)
	assets.Post("/", handlers.CreateAsset)
	assets.Patch("/:id", handlers.UpdateAsset)
	assets.Delete("/:id", handlers.DeleteAsset)

	// Ticket routes
	tickets := api.Group("/tickets")
	tickets.Get("/", handlers.GetTickets)
	tickets.Get("/:id", handlers.GetTicket)
	tickets.Post("/", handlers.CreateTicket)
	tickets.Patch("/:id", handlers.UpdateTicket)
	tickets.Delete("/:id", handlers.DeleteTicket)

	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
