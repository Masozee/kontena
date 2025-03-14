package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/your-username/project-management/database"
	"github.com/your-username/project-management/models"
)

// GetProjects returns all projects for a tenant
// @Summary Get all projects
// @Description Get all projects for the current tenant
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} models.Project
// @Router /projects [get]
func GetProjects(c *fiber.Ctx) error {
	// Get tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	var projects []models.Project
	result := database.DB.Where("tenant_id = ?", tenantID).Find(&projects)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve projects",
		})
	}

	return c.JSON(projects)
}

// GetProject returns a specific project
// @Summary Get a project
// @Description Get a project by ID for the current tenant
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [get]
func GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	var project models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&project)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	return c.JSON(project)
}

// GetProjectWithDetails returns a project with all its relationships
// @Summary Get a project with details
// @Description Get a project by ID with all its relationships for the current tenant
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project
// @Failure 404 {object} map[string]string
// @Router /projects/{id}/details [get]
func GetProjectWithDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	var project models.Project
	result := database.DB.
		Preload("People").
		Preload("KPIs").
		Preload("Tasks").
		Preload("Tasks.AssignedTo").
		Preload("Reports").
		Preload("Milestones").
		Preload("Risks").
		Preload("Issues").
		Preload("Documents").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&project)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	return c.JSON(project)
}

// CreateProject creates a new project
// @Summary Create a project
// @Description Create a new project for the current tenant
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project information"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string
// @Router /projects [post]
func CreateProject(c *fiber.Ctx) error {
	project := new(models.Project)
	if err := c.BodyParser(project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			project.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create project",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

// UpdateProject updates an existing project
// @Summary Update a project
// @Description Update an existing project for the current tenant
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Project information"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [patch]
func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	// Check if project exists
	var existingProject models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingProject)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	// Parse request body
	updatedProject := new(models.Project)
	if err := c.BodyParser(updatedProject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedProject.TenantID = existingProject.TenantID
	updatedProject.ID = existingProject.ID

	// Update project
	database.DB.Model(&existingProject).Updates(updatedProject)

	return c.JSON(existingProject)
}

// DeleteProject deletes a project
// @Summary Delete a project
// @Description Delete a project by ID for the current tenant
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /projects/{id} [delete]
func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	var project models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	database.DB.Delete(&project)

	return c.JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}
