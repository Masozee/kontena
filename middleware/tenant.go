package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// TenantMiddleware extracts tenant_id from request headers or query parameters
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get tenant_id from headers (multiple formats)
		tenantID := c.Get("X-Tenant-ID")

		// If not found, try alternate header format
		if tenantID == "" {
			tenantID = c.Get("tenant_id")
		}

		// If not in header, try query parameter
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}

		// Store tenant_id in context locals for handlers to use
		if tenantID != "" {
			// Validate that it's a valid number
			_, err := strconv.Atoi(tenantID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid tenant ID format",
				})
			}
			c.Locals("tenant_id", tenantID)
		} else {
			// If no tenant ID is provided, return an error
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tenant ID is required",
			})
		}

		return c.Next()
	}
}

// SkipTenantMiddleware is used for routes that don't require tenant isolation
func SkipTenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
