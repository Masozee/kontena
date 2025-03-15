package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/internal/database"
	"github.com/Masozee/kontena/api/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetUsers returns all users for a tenant
// @Summary Get all users
// @Description Get all users for the current tenant
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
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

	var users []models.User
	result := database.DB.Where("tenant_id = ?", tenantID).Find(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(users)
}

// GetUser returns a specific user
// @Summary Get a user
// @Description Get a user by ID for the current tenant
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUser(c *fiber.Ctx) error {
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

	var user models.User
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// CreateUser creates a new user
// @Summary Create a user
// @Description Create a new user for the current tenant
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			user.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user for the current tenant
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User information"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [patch]
func UpdateUser(c *fiber.Ctx) error {
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

	// Check if user exists
	var existingUser models.User
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingUser)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Parse request body
	updatedUser := new(models.User)
	if err := c.BodyParser(updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedUser.TenantID = existingUser.TenantID
	updatedUser.ID = existingUser.ID

	// Update user
	database.DB.Model(&existingUser).Updates(updatedUser)

	return c.JSON(existingUser)
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a user by ID for the current tenant
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
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

	var user models.User
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	database.DB.Delete(&user)

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
