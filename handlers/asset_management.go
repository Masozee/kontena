package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/models"
	"github.com/gofiber/fiber/v2"
)

// Asset Category Handlers

// GetAssetCategories returns all asset categories for a tenant
// @Summary Get all asset categories
// @Description Get all asset categories for the authenticated tenant
// @Tags asset-categories
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Success 200 {array} models.AssetCategory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /asset-categories [get]
func GetAssetCategories(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var categories []models.AssetCategory
	result := database.DB.Where("tenant_id = ?", tenantID).Find(&categories)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve asset categories",
		})
	}

	return c.JSON(categories)
}

// GetAssetCategory returns a specific asset category
// @Summary Get an asset category
// @Description Get an asset category by ID
// @Tags asset-categories
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset Category ID"
// @Success 200 {object} models.AssetCategory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /asset-categories/{id} [get]
func GetAssetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var category models.AssetCategory
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset category not found",
		})
	}

	return c.JSON(category)
}

// CreateAssetCategory creates a new asset category
// @Summary Create an asset category
// @Description Create a new asset category
// @Tags asset-categories
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param category body models.AssetCategory true "Asset Category object"
// @Success 201 {object} models.AssetCategory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /asset-categories [post]
func CreateAssetCategory(c *fiber.Ctx) error {
	category := new(models.AssetCategory)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	category.TenantID = uint(tenantID)

	// Validate required fields
	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	result := database.DB.Create(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create asset category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateAssetCategory updates an asset category
// @Summary Update an asset category
// @Description Update an asset category by ID
// @Tags asset-categories
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset Category ID"
// @Param category body models.AssetCategory true "Asset Category object"
// @Success 200 {object} models.AssetCategory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /asset-categories/{id} [put]
func UpdateAssetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var category models.AssetCategory
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset category not found",
		})
	}

	// Only update allowed fields
	updateData := new(models.AssetCategory)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Name != "" {
		category.Name = updateData.Name
	}
	category.Description = updateData.Description
	category.ParentID = updateData.ParentID

	database.DB.Save(&category)
	return c.JSON(category)
}

// DeleteAssetCategory deletes an asset category
// @Summary Delete an asset category
// @Description Delete an asset category by ID
// @Tags asset-categories
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset Category ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /asset-categories/{id} [delete]
func DeleteAssetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var category models.AssetCategory
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset category not found",
		})
	}

	// Check if category has assets
	var count int64
	database.DB.Model(&models.Asset{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete category with associated assets",
		})
	}

	database.DB.Delete(&category)
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Asset category deleted successfully",
	})
}

// Asset Handlers

// GetAssets returns all assets for a tenant
// @Summary Get all assets
// @Description Get all assets for the authenticated tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param category_id query int false "Filter by category ID"
// @Param status query string false "Filter by status"
// @Param location_id query int false "Filter by location ID"
// @Param assigned_to query int false "Filter by assignee ID"
// @Success 200 {array} models.Asset
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets [get]
func GetAssets(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Build query with filters
	query := database.DB.Where("tenant_id = ?", tenantID)

	// Apply optional filters
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if locationID := c.Query("location_id"); locationID != "" {
		query = query.Where("location_id = ?", locationID)
	}

	if assigneeID := c.Query("assigned_to"); assigneeID != "" {
		query = query.Where("current_assignee = ?", assigneeID)
	}

	// Execute query
	var assets []models.Asset
	result := query.Preload("Category").Find(&assets)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve assets",
		})
	}

	return c.JSON(assets)
}

// GetAsset returns a specific asset
// @Summary Get an asset
// @Description Get an asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset ID"
// @Success 200 {object} models.Asset
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [get]
func GetAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Category").
		Preload("Location").
		Preload("AssignedTo").
		First(&asset)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	return c.JSON(asset)
}

// CreateAsset creates a new asset
// @Summary Create an asset
// @Description Create a new asset
// @Tags assets
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param asset body models.Asset true "Asset object"
// @Success 201 {object} models.Asset
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets [post]
func CreateAsset(c *fiber.Ctx) error {
	asset := new(models.Asset)
	if err := c.BodyParser(asset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	asset.TenantID = uint(tenantID)

	// Validate required fields
	if asset.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	if asset.CategoryID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}

	// Verify category exists and belongs to tenant
	var category models.AssetCategory
	result := database.DB.Where("id = ? AND tenant_id = ?", asset.CategoryID, tenantID).First(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Set default status if not provided
	if asset.Status == "" {
		asset.Status = models.AssetStatusInStock
	}

	result = database.DB.Create(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create asset",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(asset)
}

// UpdateAsset updates an asset
// @Summary Update an asset
// @Description Update an asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset ID"
// @Param asset body models.Asset true "Asset object"
// @Success 200 {object} models.Asset
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [put]
func UpdateAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	// Only update allowed fields
	updateData := new(models.Asset)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Name != "" {
		asset.Name = updateData.Name
	}
	asset.Description = updateData.Description

	if updateData.CategoryID != 0 {
		// Verify category exists and belongs to tenant
		var category models.AssetCategory
		result := database.DB.Where("id = ? AND tenant_id = ?", updateData.CategoryID, tenantID).First(&category)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		asset.CategoryID = updateData.CategoryID
	}

	asset.SerialNumber = updateData.SerialNumber
	asset.ModelNumber = updateData.ModelNumber
	asset.Manufacturer = updateData.Manufacturer
	asset.PurchaseDate = updateData.PurchaseDate
	asset.PurchasePrice = updateData.PurchasePrice
	asset.WarrantyExpiry = updateData.WarrantyExpiry

	if updateData.Status != "" {
		asset.Status = updateData.Status
	}

	asset.LocationID = updateData.LocationID
	asset.CurrentAssignee = updateData.CurrentAssignee
	asset.ExpectedLifespan = updateData.ExpectedLifespan
	asset.Notes = updateData.Notes
	asset.Tags = updateData.Tags
	asset.Barcode = updateData.Barcode

	database.DB.Save(&asset)
	return c.JSON(asset)
}

// DeleteAsset deletes an asset
// @Summary Delete an asset
// @Description Delete an asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [delete]
func DeleteAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset not found",
		})
	}

	// Check if asset is currently assigned
	if asset.CurrentAssignee != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete asset that is currently assigned",
		})
	}

	// Check if asset has maintenance records
	var maintenanceCount int64
	database.DB.Model(&models.MaintenanceRecord{}).Where("asset_id = ?", id).Count(&maintenanceCount)
	if maintenanceCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete asset with maintenance records",
		})
	}

	database.DB.Delete(&asset)
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Asset deleted successfully",
	})
}

// Location Handlers

// GetLocations returns all locations for a tenant
// @Summary Get all locations
// @Description Get all locations for the authenticated tenant
// @Tags locations
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param type query string false "Filter by location type"
// @Success 200 {array} models.Location
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /locations [get]
func GetLocations(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Build query with filters
	query := database.DB.Where("tenant_id = ?", tenantID)

	// Apply optional filters
	if locationType := c.Query("type"); locationType != "" {
		query = query.Where("type = ?", locationType)
	}

	// Execute query
	var locations []models.Location
	result := query.Find(&locations)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve locations",
		})
	}

	return c.JSON(locations)
}

// GetLocation returns a specific location
// @Summary Get a location
// @Description Get a location by ID
// @Tags locations
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Location ID"
// @Success 200 {object} models.Location
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /locations/{id} [get]
func GetLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var location models.Location
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&location)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	return c.JSON(location)
}

// CreateLocation creates a new location
// @Summary Create a location
// @Description Create a new location
// @Tags locations
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param location body models.Location true "Location object"
// @Success 201 {object} models.Location
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /locations [post]
func CreateLocation(c *fiber.Ctx) error {
	location := new(models.Location)
	if err := c.BodyParser(location); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	location.TenantID = uint(tenantID)

	// Validate required fields
	if location.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	// If parent ID is provided, verify it exists and belongs to tenant
	if location.ParentID != nil {
		var parentLocation models.Location
		result := database.DB.Where("id = ? AND tenant_id = ?", *location.ParentID, tenantID).First(&parentLocation)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parent location ID",
			})
		}
	}

	result := database.DB.Create(&location)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create location",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(location)
}

// UpdateLocation updates a location
// @Summary Update a location
// @Description Update a location by ID
// @Tags locations
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Location ID"
// @Param location body models.Location true "Location object"
// @Success 200 {object} models.Location
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /locations/{id} [put]
func UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var location models.Location
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&location)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Only update allowed fields
	updateData := new(models.Location)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Name != "" {
		location.Name = updateData.Name
	}
	location.Description = updateData.Description
	location.Address = updateData.Address
	location.Type = updateData.Type

	// If parent ID is updated, verify it exists and belongs to tenant
	if updateData.ParentID != nil {
		// Prevent circular references
		if *updateData.ParentID == location.ID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Location cannot be its own parent",
			})
		}

		var parentLocation models.Location
		result := database.DB.Where("id = ? AND tenant_id = ?", *updateData.ParentID, tenantID).First(&parentLocation)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parent location ID",
			})
		}
		location.ParentID = updateData.ParentID
	}

	database.DB.Save(&location)
	return c.JSON(location)
}

// DeleteLocation deletes a location
// @Summary Delete a location
// @Description Delete a location by ID
// @Tags locations
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Location ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /locations/{id} [delete]
func DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var location models.Location
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&location)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if location has child locations
	var childCount int64
	database.DB.Model(&models.Location{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete location with child locations",
		})
	}

	// Check if location has assets
	var assetCount int64
	database.DB.Model(&models.Asset{}).Where("location_id = ?", id).Count(&assetCount)
	if assetCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete location with associated assets",
		})
	}

	database.DB.Delete(&location)
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Location deleted successfully",
	})
}

// Vendor Handlers

// GetVendors returns all vendors for a tenant
// @Summary Get all vendors
// @Description Get all vendors for the authenticated tenant
// @Tags vendors
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param search query string false "Search by name or contact"
// @Success 200 {array} models.Vendor
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vendors [get]
func GetVendors(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Build query with filters
	query := database.DB.Where("tenant_id = ?", tenantID)

	// Apply optional search
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR contact_name LIKE ? OR contact_email LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Execute query
	var vendors []models.Vendor
	result := query.Find(&vendors)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve vendors",
		})
	}

	return c.JSON(vendors)
}

// GetVendor returns a specific vendor
// @Summary Get a vendor
// @Description Get a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Vendor ID"
// @Success 200 {object} models.Vendor
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /vendors/{id} [get]
func GetVendor(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var vendor models.Vendor
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&vendor)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vendor not found",
		})
	}

	return c.JSON(vendor)
}

// CreateVendor creates a new vendor
// @Summary Create a vendor
// @Description Create a new vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param vendor body models.Vendor true "Vendor object"
// @Success 201 {object} models.Vendor
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vendors [post]
func CreateVendor(c *fiber.Ctx) error {
	vendor := new(models.Vendor)
	if err := c.BodyParser(vendor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	vendor.TenantID = uint(tenantID)

	// Validate required fields
	if vendor.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	result := database.DB.Create(&vendor)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create vendor",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(vendor)
}

// UpdateVendor updates a vendor
// @Summary Update a vendor
// @Description Update a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Vendor ID"
// @Param vendor body models.Vendor true "Vendor object"
// @Success 200 {object} models.Vendor
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /vendors/{id} [put]
func UpdateVendor(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var vendor models.Vendor
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&vendor)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vendor not found",
		})
	}

	// Only update allowed fields
	updateData := new(models.Vendor)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Name != "" {
		vendor.Name = updateData.Name
	}
	vendor.ContactName = updateData.ContactName
	vendor.ContactEmail = updateData.ContactEmail
	vendor.ContactPhone = updateData.ContactPhone
	vendor.Address = updateData.Address
	vendor.Website = updateData.Website
	vendor.Notes = updateData.Notes

	database.DB.Save(&vendor)
	return c.JSON(vendor)
}

// DeleteVendor deletes a vendor
// @Summary Delete a vendor
// @Description Delete a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Vendor ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /vendors/{id} [delete]
func DeleteVendor(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var vendor models.Vendor
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&vendor)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vendor not found",
		})
	}

	// Check if vendor has purchase orders
	var poCount int64
	database.DB.Model(&models.PurchaseOrder{}).Where("vendor_id = ?", id).Count(&poCount)
	if poCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete vendor with associated purchase orders",
		})
	}

	database.DB.Delete(&vendor)
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Vendor deleted successfully",
	})
}

// Procurement Request Handlers

// GetProcurementRequests returns all procurement requests for a tenant
// @Summary Get all procurement requests
// @Description Get all procurement requests for the authenticated tenant
// @Tags procurement-requests
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param status query string false "Filter by status"
// @Param requested_by query int false "Filter by requester ID"
// @Success 200 {array} models.ProcurementRequest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /procurement-requests [get]
func GetProcurementRequests(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Build query with filters
	query := database.DB.Where("tenant_id = ?", tenantID)

	// Apply optional filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if requesterID := c.Query("requested_by"); requesterID != "" {
		query = query.Where("requested_by_id = ?", requesterID)
	}

	// Execute query
	var requests []models.ProcurementRequest
	result := query.Preload("RequestedBy").Preload("ApprovedBy").Find(&requests)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve procurement requests",
		})
	}

	return c.JSON(requests)
}

// GetProcurementRequest returns a specific procurement request
// @Summary Get a procurement request
// @Description Get a procurement request by ID
// @Tags procurement-requests
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Procurement Request ID"
// @Success 200 {object} models.ProcurementRequest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /procurement-requests/{id} [get]
func GetProcurementRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var request models.ProcurementRequest
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("RequestedBy").
		Preload("ApprovedBy").
		Preload("Items").
		Preload("Items.Category").
		Preload("Items.PreferredVendor").
		Preload("PurchaseOrders").
		First(&request)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Procurement request not found",
		})
	}

	return c.JSON(request)
}

// CreateProcurementRequest creates a new procurement request
// @Summary Create a procurement request
// @Description Create a new procurement request
// @Tags procurement-requests
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param request body models.ProcurementRequest true "Procurement Request object"
// @Success 201 {object} models.ProcurementRequest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /procurement-requests [post]
func CreateProcurementRequest(c *fiber.Ctx) error {
	request := new(models.ProcurementRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	request.TenantID = uint(tenantID)

	// Validate required fields
	if request.RequestedByID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Requested by ID is required",
		})
	}

	// Verify requester exists and belongs to tenant
	var requester models.Person
	result := database.DB.Where("id = ? AND tenant_id = ?", request.RequestedByID, tenantID).First(&requester)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid requester ID",
		})
	}

	// Set default values
	if request.Status == "" {
		request.Status = models.ProcurementStatusDraft
	}

	if request.RequestNumber == "" {
		// Generate a request number (PR-YYYYMMDD-XXX)
		var lastRequest models.ProcurementRequest
		database.DB.Where("tenant_id = ?", tenantID).Order("id desc").First(&lastRequest)

		today := time.Now().Format("20060102")
		sequence := 1
		if lastRequest.ID > 0 {
			sequence = int(lastRequest.ID) + 1
		}

		request.RequestNumber = fmt.Sprintf("PR-%s-%03d", today, sequence)
	}

	if request.RequestDate.IsZero() {
		request.RequestDate = time.Now()
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Create the procurement request
	if err := tx.Create(&request).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create procurement request",
		})
	}

	// Process items if provided
	if len(request.Items) > 0 {
		for i := range request.Items {
			request.Items[i].ProcurementID = request.ID
			request.Items[i].Status = "pending"

			// Verify category exists and belongs to tenant
			var category models.AssetCategory
			result := tx.Where("id = ? AND tenant_id = ?", request.Items[i].CategoryID, tenantID).First(&category)
			if result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid category ID in item " + strconv.Itoa(i+1),
				})
			}

			// Verify preferred vendor if provided
			if request.Items[i].PreferredVendorID != nil {
				var vendor models.Vendor
				result := tx.Where("id = ? AND tenant_id = ?", *request.Items[i].PreferredVendorID, tenantID).First(&vendor)
				if result.Error != nil {
					tx.Rollback()
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Invalid preferred vendor ID in item " + strconv.Itoa(i+1),
					})
				}
			}

			if err := tx.Create(&request.Items[i]).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create procurement item",
				})
			}
		}
	}

	// Commit transaction
	tx.Commit()

	// Reload the request with all relationships
	database.DB.Where("id = ?", request.ID).
		Preload("RequestedBy").
		Preload("Items").
		Preload("Items.Category").
		Preload("Items.PreferredVendor").
		First(&request)

	return c.Status(fiber.StatusCreated).JSON(request)
}

// UpdateProcurementRequest updates a procurement request
// @Summary Update a procurement request
// @Description Update a procurement request by ID
// @Tags procurement-requests
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Procurement Request ID"
// @Param request body models.ProcurementRequest true "Procurement Request object"
// @Success 200 {object} models.ProcurementRequest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /procurement-requests/{id} [put]
func UpdateProcurementRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var request models.ProcurementRequest
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&request)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Procurement request not found",
		})
	}

	// Only allow updates to draft or submitted requests
	if request.Status != models.ProcurementStatusDraft && request.Status != models.ProcurementStatusSubmitted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot update procurement request in " + string(request.Status) + " status",
		})
	}

	// Only update allowed fields
	updateData := new(models.ProcurementRequest)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Status != "" {
		// Validate status transitions
		if request.Status == models.ProcurementStatusDraft {
			if updateData.Status != models.ProcurementStatusSubmitted &&
				updateData.Status != models.ProcurementStatusCancelled {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid status transition from draft",
				})
			}
		} else if request.Status == models.ProcurementStatusSubmitted {
			if updateData.Status != models.ProcurementStatusApproved &&
				updateData.Status != models.ProcurementStatusRejected &&
				updateData.Status != models.ProcurementStatusCancelled {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid status transition from submitted",
				})
			}

			// If approving, require approval fields
			if updateData.Status == models.ProcurementStatusApproved {
				if updateData.ApprovedByID == nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Approved by ID is required when approving",
					})
				}

				// Verify approver exists and belongs to tenant
				var approver models.Person
				result := database.DB.Where("id = ? AND tenant_id = ?", *updateData.ApprovedByID, tenantID).First(&approver)
				if result.Error != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Invalid approver ID",
					})
				}

				// Set approval date
				now := time.Now()
				request.ApprovalDate = &now
			}
		}

		request.Status = updateData.Status
	}

	request.Notes = updateData.Notes
	request.ExpectedDate = updateData.ExpectedDate
	request.TotalBudget = updateData.TotalBudget

	if updateData.ApprovedByID != nil {
		request.ApprovedByID = updateData.ApprovedByID
	}

	database.DB.Save(&request)

	// Reload the request with all relationships
	database.DB.Where("id = ?", request.ID).
		Preload("RequestedBy").
		Preload("ApprovedBy").
		Preload("Items").
		Preload("Items.Category").
		Preload("Items.PreferredVendor").
		First(&request)

	return c.JSON(request)
}

// DeleteProcurementRequest deletes a procurement request
// @Summary Delete a procurement request
// @Description Delete a procurement request by ID
// @Tags procurement-requests
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Procurement Request ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /procurement-requests/{id} [delete]
func DeleteProcurementRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var request models.ProcurementRequest
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&request)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Procurement request not found",
		})
	}

	// Only allow deletion of draft requests
	if request.Status != models.ProcurementStatusDraft {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Can only delete procurement requests in draft status",
		})
	}

	// Check if request has purchase orders
	var poCount int64
	database.DB.Model(&models.PurchaseOrder{}).Where("procurement_id = ?", id).Count(&poCount)
	if poCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete procurement request with associated purchase orders",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Delete items first
	if err := tx.Where("procurement_id = ?", id).Delete(&models.ProcurementItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete procurement items",
		})
	}

	// Delete the request
	if err := tx.Delete(&request).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete procurement request",
		})
	}

	// Commit transaction
	tx.Commit()

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Procurement request deleted successfully",
	})
}

// Asset Assignment Handlers

// GetAssetAssignments returns all asset assignments for a tenant
// @Summary Get all asset assignments
// @Description Get all asset assignments for the authenticated tenant
// @Tags asset-assignments
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param asset_id query int false "Filter by asset ID"
// @Param assigned_to query int false "Filter by assignee ID"
// @Param status query string false "Filter by status"
// @Success 200 {array} models.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /asset-assignments [get]
func GetAssetAssignments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Build query with filters
	query := database.DB.Where("tenant_id = ?", tenantID)

	// Apply optional filters
	if assetID := c.Query("asset_id"); assetID != "" {
		query = query.Where("asset_id = ?", assetID)
	}

	if assigneeID := c.Query("assigned_to"); assigneeID != "" {
		query = query.Where("assigned_to_id = ?", assigneeID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Execute query
	var assignments []models.AssetAssignment
	result := query.Preload("Asset").Preload("AssignedTo").Preload("AssignedBy").Find(&assignments)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve asset assignments",
		})
	}

	return c.JSON(assignments)
}

// GetAssetAssignment returns a specific asset assignment
// @Summary Get an asset assignment
// @Description Get an asset assignment by ID
// @Tags asset-assignments
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset Assignment ID"
// @Success 200 {object} models.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /asset-assignments/{id} [get]
func GetAssetAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var assignment models.AssetAssignment
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Asset").
		Preload("AssignedTo").
		Preload("AssignedBy").
		First(&assignment)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset assignment not found",
		})
	}

	return c.JSON(assignment)
}

// CreateAssetAssignment creates a new asset assignment
// @Summary Create an asset assignment
// @Description Create a new asset assignment
// @Tags asset-assignments
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param assignment body models.AssetAssignment true "Asset Assignment object"
// @Success 201 {object} models.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /asset-assignments [post]
func CreateAssetAssignment(c *fiber.Ctx) error {
	assignment := new(models.AssetAssignment)
	if err := c.BodyParser(assignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	assignment.TenantID = uint(tenantID)

	// Validate required fields
	if assignment.AssetID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset ID is required",
		})
	}

	if assignment.AssignedToID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assigned to ID is required",
		})
	}

	if assignment.AssignedByID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assigned by ID is required",
		})
	}

	// Verify asset exists and belongs to tenant
	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", assignment.AssetID, tenantID).First(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	// Check if asset is available for assignment
	if asset.Status != models.AssetStatusInStock && asset.Status != models.AssetStatusAssigned {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset is not available for assignment (current status: " + string(asset.Status) + ")",
		})
	}

	// Check if asset is already assigned to someone else
	if asset.Status == models.AssetStatusAssigned && asset.CurrentAssignee != nil && *asset.CurrentAssignee != assignment.AssignedToID {
		// Check if there's an active assignment for this asset
		var activeAssignment models.AssetAssignment
		result := database.DB.Where("asset_id = ? AND status = 'active' AND return_date IS NULL", assignment.AssetID).First(&activeAssignment)
		if result.Error == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Asset is already assigned to someone else",
			})
		}
	}

	// Verify assignee exists and belongs to tenant
	var assignee models.Person
	result = database.DB.Where("id = ? AND tenant_id = ?", assignment.AssignedToID, tenantID).First(&assignee)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid assignee ID",
		})
	}

	// Verify assigner exists and belongs to tenant
	var assigner models.Person
	result = database.DB.Where("id = ? AND tenant_id = ?", assignment.AssignedByID, tenantID).First(&assigner)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid assigner ID",
		})
	}

	// Set default values
	if assignment.AssignmentDate.IsZero() {
		assignment.AssignmentDate = time.Now()
	}

	if assignment.Status == "" {
		assignment.Status = "active"
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Create the assignment
	if err := tx.Create(&assignment).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create asset assignment",
		})
	}

	// Update the asset status and assignee
	asset.Status = models.AssetStatusAssigned
	asset.CurrentAssignee = &assignment.AssignedToID

	if err := tx.Save(&asset).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update asset status",
		})
	}

	// Commit transaction
	tx.Commit()

	// Reload the assignment with all relationships
	database.DB.Where("id = ?", assignment.ID).
		Preload("Asset").
		Preload("AssignedTo").
		Preload("AssignedBy").
		First(&assignment)

	return c.Status(fiber.StatusCreated).JSON(assignment)
}

// UpdateAssetAssignment updates an asset assignment
// @Summary Update an asset assignment
// @Description Update an asset assignment by ID
// @Tags asset-assignments
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset Assignment ID"
// @Param assignment body models.AssetAssignment true "Asset Assignment object"
// @Success 200 {object} models.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /asset-assignments/{id} [put]
func UpdateAssetAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var assignment models.AssetAssignment
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&assignment)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset assignment not found",
		})
	}

	// Only update allowed fields
	updateData := new(models.AssetAssignment)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Handle status changes
	if updateData.Status != "" && updateData.Status != assignment.Status {
		// If changing to returned, set return date if not provided
		if updateData.Status == "returned" && (updateData.ReturnDate == nil || updateData.ReturnDate.IsZero()) {
			now := time.Now()
			assignment.ReturnDate = &now
		}

		// If returning the asset, update the asset status
		if updateData.Status == "returned" {
			var asset models.Asset
			result := tx.Where("id = ?", assignment.AssetID).First(&asset)
			if result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to retrieve asset",
				})
			}

			// Update asset status and clear assignee
			asset.Status = models.AssetStatusInStock
			asset.CurrentAssignee = nil

			if err := tx.Save(&asset).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update asset status",
				})
			}
		}

		assignment.Status = updateData.Status
	}

	// Update other fields
	if updateData.ReturnDate != nil {
		assignment.ReturnDate = updateData.ReturnDate
	}

	if updateData.ExpectedReturn != nil {
		assignment.ExpectedReturn = updateData.ExpectedReturn
	}

	assignment.Notes = updateData.Notes

	if err := tx.Save(&assignment).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update asset assignment",
		})
	}

	// Commit transaction
	tx.Commit()

	// Reload the assignment with all relationships
	database.DB.Where("id = ?", assignment.ID).
		Preload("Asset").
		Preload("AssignedTo").
		Preload("AssignedBy").
		First(&assignment)

	return c.JSON(assignment)
}

// DeleteAssetAssignment deletes an asset assignment
// @Summary Delete an asset assignment
// @Description Delete an asset assignment by ID
// @Tags asset-assignments
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Asset Assignment ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /asset-assignments/{id} [delete]
func DeleteAssetAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var assignment models.AssetAssignment
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&assignment)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Asset assignment not found",
		})
	}

	// Only allow deletion of active assignments
	if assignment.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Can only delete active assignments",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()

	// If this is the current assignment for the asset, update the asset
	var asset models.Asset
	result = tx.Where("id = ? AND current_assignee = ?", assignment.AssetID, assignment.AssignedToID).First(&asset)
	if result.Error == nil {
		// Update asset status and clear assignee
		asset.Status = models.AssetStatusInStock
		asset.CurrentAssignee = nil

		if err := tx.Save(&asset).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update asset status",
			})
		}
	}

	// Delete the assignment
	if err := tx.Delete(&assignment).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete asset assignment",
		})
	}

	// Commit transaction
	tx.Commit()

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Asset assignment deleted successfully",
	})
}

// Maintenance Record Handlers

// GetMaintenanceRecords returns all maintenance records for a tenant
// @Summary Get all maintenance records
// @Description Get all maintenance records for the authenticated tenant
// @Tags maintenance-records
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param asset_id query int false "Filter by asset ID"
// @Param status query string false "Filter by status"
// @Param type query string false "Filter by maintenance type"
// @Success 200 {array} models.MaintenanceRecord
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance-records [get]
func GetMaintenanceRecords(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Build query with filters
	query := database.DB.Where("tenant_id = ?", tenantID)

	// Apply optional filters
	if assetID := c.Query("asset_id"); assetID != "" {
		query = query.Where("asset_id = ?", assetID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if maintenanceType := c.Query("type"); maintenanceType != "" {
		query = query.Where("maintenance_type = ?", maintenanceType)
	}

	// Execute query
	var records []models.MaintenanceRecord
	result := query.Preload("Asset").Preload("PerformedBy").Preload("Vendor").Find(&records)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve maintenance records",
		})
	}

	return c.JSON(records)
}

// GetMaintenanceRecord returns a specific maintenance record
// @Summary Get a maintenance record
// @Description Get a maintenance record by ID
// @Tags maintenance-records
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Maintenance Record ID"
// @Success 200 {object} models.MaintenanceRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /maintenance-records/{id} [get]
func GetMaintenanceRecord(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var record models.MaintenanceRecord
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Asset").
		Preload("PerformedBy").
		Preload("Vendor").
		First(&record)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Maintenance record not found",
		})
	}

	return c.JSON(record)
}

// CreateMaintenanceRecord creates a new maintenance record
// @Summary Create a maintenance record
// @Description Create a new maintenance record
// @Tags maintenance-records
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param record body models.MaintenanceRecord true "Maintenance Record object"
// @Success 201 {object} models.MaintenanceRecord
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance-records [post]
func CreateMaintenanceRecord(c *fiber.Ctx) error {
	record := new(models.MaintenanceRecord)
	if err := c.BodyParser(record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set tenant ID from context
	tenantIDStr := c.Locals("tenant_id").(string)
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}
	record.TenantID = uint(tenantID)

	// Validate required fields
	if record.AssetID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Asset ID is required",
		})
	}

	if record.MaintenanceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Maintenance type is required",
		})
	}

	if record.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is required",
		})
	}

	// Verify asset exists and belongs to tenant
	var asset models.Asset
	result := database.DB.Where("id = ? AND tenant_id = ?", record.AssetID, tenantID).First(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	// Verify performer if provided
	if record.PerformedByID != nil {
		var performer models.Person
		result := database.DB.Where("id = ? AND tenant_id = ?", *record.PerformedByID, tenantID).First(&performer)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid performer ID",
			})
		}
	}

	// Verify vendor if provided
	if record.VendorID != nil {
		var vendor models.Vendor
		result := database.DB.Where("id = ? AND tenant_id = ?", *record.VendorID, tenantID).First(&vendor)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid vendor ID",
			})
		}
	}

	// Set default values
	if record.Status == "" {
		record.Status = models.MaintenanceStatusScheduled
	}

	if record.ScheduledDate.IsZero() {
		record.ScheduledDate = time.Now()
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Create the maintenance record
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create maintenance record",
		})
	}

	// If status is in_progress, update asset status
	if record.Status == models.MaintenanceStatusInProgress {
		asset.Status = models.AssetStatusMaintenance
		if err := tx.Save(&asset).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update asset status",
			})
		}
	}

	// Commit transaction
	tx.Commit()

	// Reload the record with all relationships
	database.DB.Where("id = ?", record.ID).
		Preload("Asset").
		Preload("PerformedBy").
		Preload("Vendor").
		First(&record)

	return c.Status(fiber.StatusCreated).JSON(record)
}

// UpdateMaintenanceRecord updates a maintenance record
// @Summary Update a maintenance record
// @Description Update a maintenance record by ID
// @Tags maintenance-records
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Maintenance Record ID"
// @Param record body models.MaintenanceRecord true "Maintenance Record object"
// @Success 200 {object} models.MaintenanceRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /maintenance-records/{id} [put]
func UpdateMaintenanceRecord(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var record models.MaintenanceRecord
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&record)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Maintenance record not found",
		})
	}

	// Get the asset
	var asset models.Asset
	result = database.DB.Where("id = ?", record.AssetID).First(&asset)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve asset",
		})
	}

	// Only update allowed fields
	updateData := new(models.MaintenanceRecord)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Handle status changes
	if updateData.Status != "" && updateData.Status != record.Status {
		// If changing to in_progress, update asset status
		if updateData.Status == models.MaintenanceStatusInProgress && record.Status != models.MaintenanceStatusInProgress {
			asset.Status = models.AssetStatusMaintenance
			if err := tx.Save(&asset).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update asset status",
				})
			}
		}

		// If changing to completed, update asset status and set completed date
		if updateData.Status == models.MaintenanceStatusCompleted {
			// Set completed date if not provided
			if updateData.CompletedDate == nil || updateData.CompletedDate.IsZero() {
				now := time.Now()
				record.CompletedDate = &now
			} else {
				record.CompletedDate = updateData.CompletedDate
			}

			// Update asset status back to in_stock if it was in maintenance
			if asset.Status == models.AssetStatusMaintenance {
				asset.Status = models.AssetStatusInStock
				if err := tx.Save(&asset).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to update asset status",
					})
				}
			}
		}

		record.Status = updateData.Status
	}

	// Update other fields
	if updateData.MaintenanceType != "" {
		record.MaintenanceType = updateData.MaintenanceType
	}

	if updateData.ScheduledDate != record.ScheduledDate {
		record.ScheduledDate = updateData.ScheduledDate
	}

	if updateData.PerformedByID != nil {
		// Verify performer exists and belongs to tenant
		var performer models.Person
		result := tx.Where("id = ? AND tenant_id = ?", *updateData.PerformedByID, tenantID).First(&performer)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid performer ID",
			})
		}
		record.PerformedByID = updateData.PerformedByID
	}

	if updateData.VendorID != nil {
		// Verify vendor exists and belongs to tenant
		var vendor models.Vendor
		result := tx.Where("id = ? AND tenant_id = ?", *updateData.VendorID, tenantID).First(&vendor)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid vendor ID",
			})
		}
		record.VendorID = updateData.VendorID
	}

	record.Cost = updateData.Cost
	record.Description = updateData.Description
	record.Results = updateData.Results
	record.NextScheduled = updateData.NextScheduled

	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update maintenance record",
		})
	}

	// Commit transaction
	tx.Commit()

	// Reload the record with all relationships
	database.DB.Where("id = ?", record.ID).
		Preload("Asset").
		Preload("PerformedBy").
		Preload("Vendor").
		First(&record)

	return c.JSON(record)
}

// DeleteMaintenanceRecord deletes a maintenance record
// @Summary Delete a maintenance record
// @Description Delete a maintenance record by ID
// @Tags maintenance-records
// @Accept json
// @Produce json
// @Param tenant_id query string false "Tenant ID"
// @Param id path int true "Maintenance Record ID"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /maintenance-records/{id} [delete]
func DeleteMaintenanceRecord(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var record models.MaintenanceRecord
	result := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&record)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Maintenance record not found",
		})
	}

	// Only allow deletion of scheduled or cancelled records
	if record.Status != models.MaintenanceStatusScheduled && record.Status != models.MaintenanceStatusCancelled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Can only delete scheduled or cancelled maintenance records",
		})
	}

	if err := database.DB.Delete(&record).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete maintenance record",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Maintenance record deleted successfully",
	})
}
