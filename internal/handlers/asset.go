package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/internal/database"
	"github.com/Masozee/kontena/api/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetAssets returns all assets for a tenant
// @Summary Get all assets
// @Description Get all assets for the current tenant
// @Tags assets
// @Accept json
// @Produce json
// @Success 200 {array} models.Asset
// @Router /assets [get]
func GetAssets(c *fiber.Ctx) error {
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

	var assets []models.Asset
	query := database.DB.Where("tenant_id = ?", tenantID)
	result := query.Find(&assets)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve assets",
		})
	}

	return c.JSON(assets)
}

// GetAsset returns a specific asset
// @Summary Get an asset
// @Description Get an asset by ID for the current tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} models.Asset
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [get]
func GetAsset(c *fiber.Ctx) error {
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

	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&asset)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	return c.JSON(asset)
}

// CreateAsset creates a new asset
// @Summary Create an asset
// @Description Create a new asset for the current tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param asset body models.Asset true "Asset information"
// @Success 201 {object} models.Asset
// @Failure 400 {object} map[string]string
// @Router /assets [post]
func CreateAsset(c *fiber.Ctx) error {
	asset := new(models.Asset)
	if err := c.BodyParser(asset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			asset.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create asset",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(asset)
}

// UpdateAsset updates an existing asset
// @Summary Update an asset
// @Description Update an existing asset for the current tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param asset body models.Asset true "Asset information"
// @Success 200 {object} models.Asset
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [patch]
func UpdateAsset(c *fiber.Ctx) error {
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

	// Check if asset exists
	var existingAsset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingAsset)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	// Parse request body
	updatedAsset := new(models.Asset)
	if err := c.BodyParser(updatedAsset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedAsset.TenantID = existingAsset.TenantID
	updatedAsset.ID = existingAsset.ID

	// Update asset
	database.DB.Model(&existingAsset).Updates(updatedAsset)

	return c.JSON(existingAsset)
}

// DeleteAsset deletes an asset
// @Summary Delete an asset
// @Description Delete an asset by ID for the current tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [delete]
func DeleteAsset(c *fiber.Ctx) error {
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

	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	database.DB.Delete(&asset)

	return c.JSON(fiber.Map{
		"message": "Asset deleted successfully",
	})
}
