package handlers

import (
	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/models"
	"github.com/gofiber/fiber/v2"
)

// GetTenants returns all tenants
// @Summary Get all tenants
// @Description Get all tenants in the system
// @Tags tenants
// @Accept json
// @Produce json
// @Success 200 {array} models.Tenant
// @Router /tenants [get]
func GetTenants(c *fiber.Ctx) error {
	var tenants []models.Tenant
	result := database.DB.Find(&tenants)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tenants",
		})
	}
	return c.JSON(tenants)
}

// GetTenant returns a specific tenant
// @Summary Get a tenant
// @Description Get a tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Success 200 {object} models.Tenant
// @Failure 404 {object} map[string]string
// @Router /tenants/{id} [get]
func GetTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	var tenant models.Tenant
	result := database.DB.First(&tenant, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}
	return c.JSON(tenant)
}

// CreateTenant creates a new tenant
// @Summary Create a tenant
// @Description Create a new tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body models.Tenant true "Tenant object"
// @Success 201 {object} models.Tenant
// @Failure 400 {object} map[string]string
// @Router /tenants [post]
func CreateTenant(c *fiber.Ctx) error {
	tenant := new(models.Tenant)
	if err := c.BodyParser(tenant); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	result := database.DB.Create(&tenant)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tenant",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tenant)
}

// UpdateTenant updates a tenant
// @Summary Update a tenant
// @Description Update a tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Param tenant body models.Tenant true "Tenant object"
// @Success 200 {object} models.Tenant
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tenants/{id} [put]
func UpdateTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	var tenant models.Tenant
	result := database.DB.First(&tenant, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	if err := c.BodyParser(&tenant); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	database.DB.Save(&tenant)
	return c.JSON(tenant)
}

// DeleteTenant deletes a tenant
// @Summary Delete a tenant
// @Description Delete a tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Success 204 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tenants/{id} [delete]
func DeleteTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	var tenant models.Tenant
	result := database.DB.First(&tenant, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	database.DB.Delete(&tenant)
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Tenant deleted successfully",
	})
}
