package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kontena/api/internal/database"
	"github.com/kontena/api/internal/models"
)

// GetLeads returns all leads for a tenant
// @Summary Get all leads
// @Description Get all leads for the current tenant
// @Tags leads
// @Accept json
// @Produce json
// @Success 200 {array} models.Lead
// @Router /leads [get]
func GetLeads(c *fiber.Ctx) error {
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

	var leads []models.Lead
	result := database.DB.Where("tenant_id = ?", tenantID).Find(&leads)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve leads",
		})
	}

	return c.JSON(leads)
}

// GetLead returns a specific lead
// @Summary Get a lead
// @Description Get a lead by ID for the current tenant
// @Tags leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Success 200 {object} models.Lead
// @Failure 404 {object} map[string]string
// @Router /leads/{id} [get]
func GetLead(c *fiber.Ctx) error {
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

	var lead models.Lead
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&lead)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lead not found",
		})
	}

	return c.JSON(lead)
}

// CreateLead creates a new lead
// @Summary Create a lead
// @Description Create a new lead for the current tenant
// @Tags leads
// @Accept json
// @Produce json
// @Param lead body models.Lead true "Lead information"
// @Success 201 {object} models.Lead
// @Failure 400 {object} map[string]string
// @Router /leads [post]
func CreateLead(c *fiber.Ctx) error {
	lead := new(models.Lead)
	if err := c.BodyParser(lead); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			lead.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&lead)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create lead",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(lead)
}

// UpdateLead updates an existing lead
// @Summary Update a lead
// @Description Update an existing lead for the current tenant
// @Tags leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Param lead body models.Lead true "Lead information"
// @Success 200 {object} models.Lead
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /leads/{id} [patch]
func UpdateLead(c *fiber.Ctx) error {
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

	// Check if lead exists
	var existingLead models.Lead
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingLead)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lead not found",
		})
	}

	// Parse request body
	updatedLead := new(models.Lead)
	if err := c.BodyParser(updatedLead); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedLead.TenantID = existingLead.TenantID
	updatedLead.ID = existingLead.ID

	// Update lead
	database.DB.Model(&existingLead).Updates(updatedLead)

	return c.JSON(existingLead)
}

// DeleteLead deletes a lead
// @Summary Delete a lead
// @Description Delete a lead by ID for the current tenant
// @Tags leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /leads/{id} [delete]
func DeleteLead(c *fiber.Ctx) error {
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

	var lead models.Lead
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&lead)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lead not found",
		})
	}

	database.DB.Delete(&lead)

	return c.JSON(fiber.Map{
		"message": "Lead deleted successfully",
	})
}
