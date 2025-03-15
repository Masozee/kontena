package handlers

import (
	"fmt"
	"strconv"

	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/models"
	"github.com/gofiber/fiber/v2"
)

// GetPeople retrieves all people for the current tenant
// @Summary Get all people
// @Description Get all people for the current tenant
// @Tags people
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Success 200 {array} models.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [get]
func GetPeople(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id").(string)
	if tenantIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing tenant_id",
		})
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid tenant_id format",
		})
	}

	var people []models.Person
	result := database.DB.Where("tenant_id = ?", tenantID).Find(&people)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve people",
		})
	}

	return c.JSON(people)
}

// GetPerson retrieves a specific person by ID for the current tenant
// @Summary Get a person
// @Description Get a person by ID for the current tenant
// @Tags people
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "Person ID"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [get]
func GetPerson(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id").(string)
	if tenantIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing tenant_id",
		})
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid tenant_id format",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var person models.Person
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&person)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Person not found",
		})
	}

	return c.JSON(person)
}

// CreatePerson creates a new person for the current tenant
// @Summary Create a person
// @Description Create a new person for the current tenant
// @Tags people
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param person body models.Person true "Person object"
// @Success 201 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [post]
func CreatePerson(c *fiber.Ctx) error {
	// Log the request body for debugging
	body := string(c.Body())
	fmt.Printf("Person request body: %s\n", body)

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

	person := new(models.Person)
	if err := c.BodyParser(person); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if person.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Person name is required",
		})
	}

	if person.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Person email is required",
		})
	}

	if person.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Person role is required",
		})
	}

	person.TenantID = uint(tenantID)

	result := database.DB.Create(&person)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create person: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(person)
}

// UpdatePerson updates an existing person by ID for the current tenant
// @Summary Update a person
// @Description Update an existing person by ID for the current tenant
// @Tags people
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "Person ID"
// @Param person body models.Person true "Person object"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [put]
func UpdatePerson(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id").(string)
	if tenantIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing tenant_id",
		})
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid tenant_id format",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var existingPerson models.Person
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingPerson)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Person not found",
		})
	}

	updatedPerson := new(models.Person)
	if err := c.BodyParser(updatedPerson); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID cannot be changed
	updatedPerson.TenantID = uint(tenantID)
	updatedPerson.ID = uint(id)

	result = database.DB.Model(&existingPerson).Updates(updatedPerson)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update person",
		})
	}

	return c.JSON(updatedPerson)
}

// DeletePerson deletes a person by ID for the current tenant
// @Summary Delete a person
// @Description Delete a person by ID for the current tenant
// @Tags people
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "Person ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [delete]
func DeletePerson(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id").(string)
	if tenantIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing tenant_id",
		})
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid tenant_id format",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var person models.Person
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&person)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Person not found",
		})
	}

	result = database.DB.Delete(&person)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete person",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Person deleted successfully",
	})
}
