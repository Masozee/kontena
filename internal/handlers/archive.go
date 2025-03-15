package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/internal/database"
	"github.com/Masozee/kontena/api/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetArchives returns all archives for a tenant
// @Summary Get all archives
// @Description Get all archives for the current tenant
// @Tags archives
// @Accept json
// @Produce json
// @Success 200 {array} models.Archive
// @Router /archives [get]
func GetArchives(c *fiber.Ctx) error {
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

	var archives []models.Archive
	query := database.DB.Where("tenant_id = ?", tenantID)
	result := query.Find(&archives)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve archives",
		})
	}

	return c.JSON(archives)
}

// GetArchive returns a specific archive
// @Summary Get an archive
// @Description Get an archive by ID for the current tenant
// @Tags archives
// @Accept json
// @Produce json
// @Param id path int true "Archive ID"
// @Success 200 {object} models.Archive
// @Failure 404 {object} map[string]string
// @Router /archives/{id} [get]
func GetArchive(c *fiber.Ctx) error {
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

	var archive models.Archive
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&archive)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Archive not found",
		})
	}

	return c.JSON(archive)
}

// CreateArchive creates a new archive
// @Summary Create an archive
// @Description Create a new archive for the current tenant
// @Tags archives
// @Accept json
// @Produce json
// @Param archive body models.Archive true "Archive information"
// @Success 201 {object} models.Archive
// @Failure 400 {object} map[string]string
// @Router /archives [post]
func CreateArchive(c *fiber.Ctx) error {
	archive := new(models.Archive)
	if err := c.BodyParser(archive); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			archive.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&archive)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create archive",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(archive)
}

// UpdateArchive updates an existing archive
// @Summary Update an archive
// @Description Update an existing archive for the current tenant
// @Tags archives
// @Accept json
// @Produce json
// @Param id path int true "Archive ID"
// @Param archive body models.Archive true "Archive information"
// @Success 200 {object} models.Archive
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /archives/{id} [patch]
func UpdateArchive(c *fiber.Ctx) error {
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

	// Check if archive exists
	var existingArchive models.Archive
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingArchive)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Archive not found",
		})
	}

	// Parse request body
	updatedArchive := new(models.Archive)
	if err := c.BodyParser(updatedArchive); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedArchive.TenantID = existingArchive.TenantID
	updatedArchive.ID = existingArchive.ID

	// Update archive
	database.DB.Model(&existingArchive).Updates(updatedArchive)

	return c.JSON(existingArchive)
}

// DeleteArchive deletes an archive
// @Summary Delete an archive
// @Description Delete an archive by ID for the current tenant
// @Tags archives
// @Accept json
// @Produce json
// @Param id path int true "Archive ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /archives/{id} [delete]
func DeleteArchive(c *fiber.Ctx) error {
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

	var archive models.Archive
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&archive)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Archive not found",
		})
	}

	database.DB.Delete(&archive)

	return c.JSON(fiber.Map{
		"message": "Archive deleted successfully",
	})
}
