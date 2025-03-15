package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/internal/database"
	"github.com/Masozee/kontena/api/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetTickets returns all tickets for a tenant
// @Summary Get all tickets
// @Description Get all tickets for the current tenant
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {array} models.Ticket
// @Router /tickets [get]
func GetTickets(c *fiber.Ctx) error {
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

	var tickets []models.Ticket
	query := database.DB.Where("tenant_id = ?", tenantID)
	result := query.Find(&tickets)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tickets",
		})
	}

	return c.JSON(tickets)
}

// GetTicket returns a specific ticket
// @Summary Get a ticket
// @Description Get a ticket by ID for the current tenant
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.Ticket
// @Failure 404 {object} map[string]string
// @Router /tickets/{id} [get]
func GetTicket(c *fiber.Ctx) error {
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

	var ticket models.Ticket
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&ticket)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	return c.JSON(ticket)
}

// CreateTicket creates a new ticket
// @Summary Create a ticket
// @Description Create a new ticket for the current tenant
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body models.Ticket true "Ticket information"
// @Success 201 {object} models.Ticket
// @Failure 400 {object} map[string]string
// @Router /tickets [post]
func CreateTicket(c *fiber.Ctx) error {
	ticket := new(models.Ticket)
	if err := c.BodyParser(ticket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenantID")
	if tenantIDStr != nil {
		tenantID, err := strconv.Atoi(tenantIDStr.(string))
		if err == nil {
			ticket.TenantID = uint(tenantID)
		}
	}

	result := database.DB.Create(&ticket)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create ticket",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ticket)
}

// UpdateTicket updates an existing ticket
// @Summary Update a ticket
// @Description Update an existing ticket for the current tenant
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param ticket body models.Ticket true "Ticket information"
// @Success 200 {object} models.Ticket
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/{id} [patch]
func UpdateTicket(c *fiber.Ctx) error {
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

	// Check if ticket exists
	var existingTicket models.Ticket
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&existingTicket)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	// Parse request body
	updatedTicket := new(models.Ticket)
	if err := c.BodyParser(updatedTicket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure tenant ID doesn't change
	updatedTicket.TenantID = existingTicket.TenantID
	updatedTicket.ID = existingTicket.ID

	// Update ticket
	database.DB.Model(&existingTicket).Updates(updatedTicket)

	return c.JSON(existingTicket)
}

// DeleteTicket deletes a ticket
// @Summary Delete a ticket
// @Description Delete a ticket by ID for the current tenant
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/{id} [delete]
func DeleteTicket(c *fiber.Ctx) error {
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

	var ticket models.Ticket
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&ticket)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket not found",
		})
	}

	database.DB.Delete(&ticket)

	return c.JSON(fiber.Map{
		"message": "Ticket deleted successfully",
	})
}
