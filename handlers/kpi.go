package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/models"
	"github.com/gofiber/fiber/v2"
)

// GetKPIs retrieves all KPIs for a specific project
// @Summary Get all KPIs for a project
// @Description Get all KPIs for a specific project
// @Tags kpis
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param project_id path int true "Project ID"
// @Success 200 {array} models.KPI
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{project_id}/kpis [get]
func GetKPIs(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
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

	projectID, err := c.ParamsInt("project_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	// Verify project belongs to tenant
	var project models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	var kpis []models.KPI
	result = database.DB.Where("project_id = ?", projectID).Find(&kpis)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve KPIs",
		})
	}

	return c.JSON(kpis)
}

// GetKPI retrieves a specific KPI by ID
// @Summary Get a KPI by ID
// @Description Get a specific KPI by ID
// @Tags kpis
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "KPI ID"
// @Success 200 {object} models.KPI
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kpis/{id} [get]
func GetKPI(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid KPI ID format",
		})
	}

	var kpi models.KPI
	result := database.DB.First(&kpi, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "KPI not found",
		})
	}

	// Verify KPI belongs to a project owned by the tenant
	var project models.Project
	result = database.DB.Where("id = ? AND tenant_id = ?", kpi.ProjectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this KPI",
		})
	}

	return c.JSON(kpi)
}

// CreateKPI creates a new KPI for a project
// @Summary Create a KPI
// @Description Create a new KPI for a project
// @Tags kpis
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param project_id path int true "Project ID"
// @Param kpi body models.KPI true "KPI object"
// @Success 201 {object} models.KPI
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{project_id}/kpis [post]
func CreateKPI(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
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

	projectID, err := c.ParamsInt("project_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	// Verify project belongs to tenant
	var project models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	kpi := new(models.KPI)
	if err := c.BodyParser(kpi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if kpi.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "KPI description is required",
		})
	}

	if kpi.TargetValue <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "KPI target value must be greater than 0",
		})
	}

	if kpi.Unit == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "KPI unit is required",
		})
	}

	kpi.ProjectID = uint(projectID)
	kpi.UpdateAchievement() // Set achieved status based on current value

	result = database.DB.Create(&kpi)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create KPI: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(kpi)
}

// UpdateKPI updates an existing KPI by ID
// @Summary Update a KPI
// @Description Update an existing KPI by ID
// @Tags kpis
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "KPI ID"
// @Param kpi body models.KPI true "KPI object"
// @Success 200 {object} models.KPI
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kpis/{id} [patch]
func UpdateKPI(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid KPI ID format",
		})
	}

	var existingKPI models.KPI
	result := database.DB.First(&existingKPI, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "KPI not found",
		})
	}

	// Verify KPI belongs to a project owned by the tenant
	var project models.Project
	result = database.DB.Where("id = ? AND tenant_id = ?", existingKPI.ProjectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this KPI",
		})
	}

	updatedKPI := new(models.KPI)
	if err := c.BodyParser(updatedKPI); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Ensure project ID cannot be changed
	updatedKPI.ProjectID = existingKPI.ProjectID
	updatedKPI.ID = uint(id)

	// Update achievement status based on current value
	if updatedKPI.CurrentValue > 0 || updatedKPI.TargetValue > 0 {
		updatedKPI.UpdateAchievement()
	}

	result = database.DB.Model(&existingKPI).Updates(updatedKPI)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update KPI: " + result.Error.Error(),
		})
	}

	// Get the updated KPI
	database.DB.First(&existingKPI, id)
	return c.JSON(existingKPI)
}

// DeleteKPI deletes a KPI by ID
// @Summary Delete a KPI
// @Description Delete a KPI by ID
// @Tags kpis
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "KPI ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kpis/{id} [delete]
func DeleteKPI(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid KPI ID format",
		})
	}

	var kpi models.KPI
	result := database.DB.First(&kpi, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "KPI not found",
		})
	}

	// Verify KPI belongs to a project owned by the tenant
	var project models.Project
	result = database.DB.Where("id = ? AND tenant_id = ?", kpi.ProjectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this KPI",
		})
	}

	result = database.DB.Delete(&kpi)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete KPI: " + result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "KPI deleted successfully",
	})
}
