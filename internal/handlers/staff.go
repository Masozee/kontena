package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kontena/api/internal/database"
	"github.com/kontena/api/internal/models"
)

// GetStaff returns all staff members for a tenant
// @Summary Get all staff members
// @Description Get all staff members for the current tenant
// @Tags staff
// @Accept json
// @Produce json
// @Success 200 {array} models.Staff
// @Router /staff [get]
func GetStaff(c *fiber.Ctx) error {
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

	var staff []models.Staff
	query := database.DB.Where("tenant_id = ?", tenantID)
	result := query.Find(&staff)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve staff members",
		})
	}

	return c.JSON(staff)
}

// GetStaffMember returns a specific staff member
// @Summary Get a staff member
// @Description Get a staff member by ID for the current tenant
// @Tags staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 {object} models.Staff
// @Failure 404 {object} map[string]string
// @Router /staff/{id} [get]
func GetStaffMember(c *fiber.Ctx) error {
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

	var staff models.Staff
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&staff)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff member not found",
		})
	}

	return c.JSON(staff)
}

// CreateStaffMember creates a new staff member
// @Summary Create a staff member
// @Description Create a new staff member for the current tenant
// @Tags staff
// @Accept json
// @Produce json
// @Param staff body models.Staff true "Staff member information"
// @Success 201 {object} models.Staff
// @Failure 400 {object} map[string]string
// @Router /staff [post]
func CreateStaffMember(c *fiber.Ctx) error {
	staff := new(models.Staff)
	if err := c.BodyParser(staff); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			staff.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&staff)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create staff member",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(staff)
}

// UpdateStaffMember updates an existing staff member
// @Summary Update a staff member
// @Description Update an existing staff member for the current tenant
// @Tags staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Param staff body models.Staff true "Staff member information"
// @Success 200 {object} models.Staff
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /staff/{id} [patch]
func UpdateStaffMember(c *fiber.Ctx) error {
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

	// Check if staff member exists
	var existingStaff models.Staff
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingStaff)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff member not found",
		})
	}

	// Parse request body
	updatedStaff := new(models.Staff)
	if err := c.BodyParser(updatedStaff); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedStaff.TenantID = existingStaff.TenantID
	updatedStaff.ID = existingStaff.ID

	// Update staff member
	database.DB.Model(&existingStaff).Updates(updatedStaff)

	return c.JSON(existingStaff)
}

// DeleteStaffMember deletes a staff member
// @Summary Delete a staff member
// @Description Delete a staff member by ID for the current tenant
// @Tags staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /staff/{id} [delete]
func DeleteStaffMember(c *fiber.Ctx) error {
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

	var staff models.Staff
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&staff)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff member not found",
		})
	}

	database.DB.Delete(&staff)

	return c.JSON(fiber.Map{
		"message": "Staff member deleted successfully",
	})
}
