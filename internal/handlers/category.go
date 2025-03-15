package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/internal/database"
	"github.com/Masozee/kontena/api/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetCategories returns all categories for a tenant
// @Summary Get all categories
// @Description Get all categories for the current tenant
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func GetCategories(c *fiber.Ctx) error {
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

	var categories []models.Category
	result := database.DB.Where("tenant_id = ?", tenantID).Find(&categories)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}

	return c.JSON(categories)
}

// GetCategory returns a specific category
// @Summary Get a category
// @Description Get a category by ID for the current tenant
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func GetCategory(c *fiber.Ctx) error {
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

	var category models.Category
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.JSON(category)
}

// CreateCategory creates a new category
// @Summary Create a category
// @Description Create a new category for the current tenant
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category information"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Router /categories [post]
func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			category.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateCategory updates an existing category
// @Summary Update a category
// @Description Update an existing category for the current tenant
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category information"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [patch]
func UpdateCategory(c *fiber.Ctx) error {
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

	// Check if category exists
	var existingCategory models.Category
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingCategory)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Parse request body
	updatedCategory := new(models.Category)
	if err := c.BodyParser(updatedCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedCategory.TenantID = existingCategory.TenantID
	updatedCategory.ID = existingCategory.ID

	// Update category
	database.DB.Model(&existingCategory).Updates(updatedCategory)

	return c.JSON(existingCategory)
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Delete a category by ID for the current tenant
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [delete]
func DeleteCategory(c *fiber.Ctx) error {
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

	var category models.Category
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	database.DB.Delete(&category)

	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
