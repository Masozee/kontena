package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kontena/api/internal/database"
	"github.com/kontena/api/internal/models"
)

func main() {
	// Initialize database connection
	database.InitDB()

	// Create tenants
	tenants := []models.Tenant{
		{
			Name:   "Acme Corporation",
			Plan:   "Enterprise",
			Status: "Active",
		},
		{
			Name:   "Stark Industries",
			Plan:   "Professional",
			Status: "Active",
		},
		{
			Name:   "Wayne Enterprises",
			Plan:   "Basic",
			Status: "Trial",
		},
	}

	for i := range tenants {
		result := database.DB.Create(&tenants[i])
		if result.Error != nil {
			log.Fatalf("Failed to create tenant: %v", result.Error)
		}
		fmt.Printf("Created tenant: %s (ID: %d)\n", tenants[i].Name, tenants[i].ID)
	}

	// Create users for each tenant
	userRoles := []models.UserRole{models.RoleAdmin, models.RoleSales, models.RoleSupport}

	for _, tenant := range tenants {
		for i, role := range userRoles {
			// Create a profile as JSON
			profile := map[string]interface{}{
				"department": "Sales",
				"location":   "New York",
				"skills":     []string{"Communication", "Negotiation"},
				"languages":  []string{"English", "Spanish"},
			}

			profileJSON, _ := json.Marshal(profile)

			user := models.User{
				TenantID: tenant.ID,
				Name:     fmt.Sprintf("User %d", i+1),
				Email:    fmt.Sprintf("user%d@%s.com", i+1, tenant.Name),
				Role:     role,
				Profile:  string(profileJSON),
			}

			result := database.DB.Create(&user)
			if result.Error != nil {
				log.Fatalf("Failed to create user: %v", result.Error)
			}
			fmt.Printf("Created user: %s (ID: %d) for tenant %s\n", user.Name, user.ID, tenant.Name)
		}
	}

	// Create staff members for each tenant
	staffRoles := []models.StaffRole{models.RoleStaffAdmin, models.RoleStaffManager, models.RoleStaffEmployee}
	departments := []string{"IT", "HR", "Finance", "Operations"}

	for _, tenant := range tenants {
		for i, role := range staffRoles {
			for j, dept := range departments {
				// Create a profile as JSON
				profile := map[string]interface{}{
					"skills":     []string{"Leadership", "Communication", "Problem Solving"},
					"education":  "Bachelor's Degree",
					"experience": fmt.Sprintf("%d years", (i+1)*2),
					"certifications": []string{
						"Project Management Professional",
						"Certified ScrumMaster",
					},
				}

				profileJSON, _ := json.Marshal(profile)

				staff := models.Staff{
					TenantID:   tenant.ID,
					Name:       fmt.Sprintf("Staff %d-%d", i+1, j+1),
					Email:      fmt.Sprintf("staff%d%d@%s.com", i+1, j+1, tenant.Name),
					Role:       role,
					Department: dept,
					Position:   fmt.Sprintf("%s %s", role, dept),
					Profile:    string(profileJSON),
				}

				result := database.DB.Create(&staff)
				if result.Error != nil {
					log.Fatalf("Failed to create staff: %v", result.Error)
				}
				fmt.Printf("Created staff: %s (ID: %d) for tenant %s\n", staff.Name, staff.ID, tenant.Name)
			}
		}
	}

	// Create categories for each tenant
	for _, tenant := range tenants {
		categories := []models.Category{
			{
				TenantID:    tenant.ID,
				Name:        "VIP Clients",
				Permissions: `["read", "write", "delete"]`,
			},
			{
				TenantID:    tenant.ID,
				Name:        "Regular Clients",
				Permissions: `["read", "write"]`,
			},
			{
				TenantID:    tenant.ID,
				Name:        "Prospects",
				Permissions: `["read"]`,
			},
		}

		for i := range categories {
			result := database.DB.Create(&categories[i])
			if result.Error != nil {
				log.Fatalf("Failed to create category: %v", result.Error)
			}
			fmt.Printf("Created category: %s (ID: %d) for tenant %s\n", categories[i].Name, categories[i].ID, tenant.Name)
		}
	}

	// Create leads for each tenant
	for _, tenant := range tenants {
		// Get categories for this tenant
		var categories []models.Category
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&categories)

		// Get users for this tenant
		var users []models.User
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&users)

		if len(categories) == 0 || len(users) == 0 {
			continue
		}

		leads := []models.Lead{
			{
				TenantID:   tenant.ID,
				Name:       "John Doe",
				Contact:    "123-456-7890",
				Email:      "john.doe@example.com",
				CategoryID: &categories[0].ID,
				AssignedTo: &users[0].ID,
				Status:     "New",
				Notes:      "Potential client for our new product",
			},
			{
				TenantID:   tenant.ID,
				Name:       "Jane Smith",
				Contact:    "987-654-3210",
				Email:      "jane.smith@example.com",
				CategoryID: &categories[1].ID,
				AssignedTo: &users[1].ID,
				Status:     "In Progress",
				Notes:      "Follow up on proposal",
			},
			{
				TenantID:   tenant.ID,
				Name:       "Bob Johnson",
				Contact:    "555-123-4567",
				Email:      "bob.johnson@example.com",
				CategoryID: &categories[2].ID,
				AssignedTo: &users[2].ID,
				Status:     "Closed",
				Notes:      "Deal closed successfully",
			},
		}

		for i := range leads {
			result := database.DB.Create(&leads[i])
			if result.Error != nil {
				log.Fatalf("Failed to create lead: %v", result.Error)
			}
			fmt.Printf("Created lead: %s (ID: %d) for tenant %s\n", leads[i].Name, leads[i].ID, tenant.Name)
		}
	}

	// Create archives for each tenant
	for _, tenant := range tenants {
		// Get categories for this tenant
		var categories []models.Category
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&categories)

		// Get staff for this tenant
		var staffMembers []models.Staff
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&staffMembers)

		if len(categories) == 0 || len(staffMembers) == 0 {
			continue
		}

		archiveStatuses := []models.ArchiveStatus{
			models.ArchiveStatusActive,
			models.ArchiveStatusInactive,
			models.ArchiveStatusConfidential,
		}

		for i, status := range archiveStatuses {
			// Create metadata as JSON
			metadata := map[string]interface{}{
				"tags":       []string{"important", "document", "contract"},
				"department": "Legal",
				"expiry":     "2025-12-31",
			}

			metadataJSON, _ := json.Marshal(metadata)

			categoryID := categories[i%len(categories)].ID
			createdByID := staffMembers[i%len(staffMembers)].ID

			archive := models.Archive{
				TenantID:    tenant.ID,
				Title:       fmt.Sprintf("Archive Document %d", i+1),
				Description: fmt.Sprintf("This is a sample archive document %d", i+1),
				CategoryID:  &categoryID,
				FilePath:    fmt.Sprintf("/uploads/tenant_%d/archive_%d.pdf", tenant.ID, i+1),
				FileType:    "application/pdf",
				FileSize:    int64(1024 * 1024 * (i + 1)), // 1MB, 2MB, 3MB
				Status:      status,
				CreatedByID: &createdByID,
				Metadata:    string(metadataJSON),
			}

			result := database.DB.Create(&archive)
			if result.Error != nil {
				log.Fatalf("Failed to create archive: %v", result.Error)
			}
			fmt.Printf("Created archive: %s (ID: %d) for tenant %s\n", archive.Title, archive.ID, tenant.Name)
		}
	}

	// Create assets for each tenant
	for _, tenant := range tenants {
		// Get staff for this tenant
		var staffMembers []models.Staff
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&staffMembers)

		if len(staffMembers) == 0 {
			continue
		}

		assetTypes := []models.AssetType{
			models.AssetTypeComputer,
			models.AssetTypeVehicle,
			models.AssetTypeFurniture,
			models.AssetTypeEquipment,
		}

		assetStatuses := []models.AssetStatus{
			models.AssetStatusAvailable,
			models.AssetStatusAssigned,
			models.AssetStatusMaintenance,
		}

		for i, assetType := range assetTypes {
			for j, status := range assetStatuses {
				// Create metadata as JSON
				metadata := map[string]interface{}{
					"manufacturer": "Example Corp",
					"warranty":     "2 years",
					"color":        "Black",
					"dimensions":   "15x10x2 inches",
				}

				metadataJSON, _ := json.Marshal(metadata)

				asset := models.Asset{
					TenantID:      tenant.ID,
					Name:          fmt.Sprintf("%s %d", assetType, i+1),
					Description:   fmt.Sprintf("This is a sample %s asset", assetType),
					Type:          assetType,
					SerialNumber:  fmt.Sprintf("SN-%d-%d-%d", tenant.ID, i+1, j+1),
					PurchaseDate:  time.Now().AddDate(0, -(i + j), 0), // Purchased i+j months ago
					PurchasePrice: float64((i + 1) * (j + 1) * 100),   // $100, $200, $300, etc.
					Status:        status,
					Location:      fmt.Sprintf("Office %d", j+1),
					Notes:         fmt.Sprintf("Sample notes for %s %d", assetType, i+1),
					Metadata:      string(metadataJSON),
				}

				// Only set AssignedToID if status is "assigned"
				if status == models.AssetStatusAssigned && len(staffMembers) > 0 {
					staffID := staffMembers[j%len(staffMembers)].ID
					asset.AssignedToID = &staffID
				}

				result := database.DB.Create(&asset)
				if result.Error != nil {
					log.Fatalf("Failed to create asset: %v", result.Error)
				}
				fmt.Printf("Created asset: %s (ID: %d) for tenant %s\n", asset.Name, asset.ID, tenant.Name)
			}
		}
	}

	// Create tickets for each tenant
	for _, tenant := range tenants {
		// Get staff for this tenant
		var staffMembers []models.Staff
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&staffMembers)

		// Get categories for this tenant
		var categories []models.Category
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&categories)

		// Get assets for this tenant
		var assets []models.Asset
		database.DB.Where("tenant_id = ?", tenant.ID).Find(&assets)

		if len(staffMembers) == 0 || len(categories) == 0 || len(assets) == 0 {
			continue
		}

		ticketPriorities := []models.TicketPriority{
			models.TicketPriorityLow,
			models.TicketPriorityMedium,
			models.TicketPriorityHigh,
			models.TicketPriorityCritical,
		}

		ticketStatuses := []models.TicketStatus{
			models.TicketStatusOpen,
			models.TicketStatusInProgress,
			models.TicketStatusResolved,
			models.TicketStatusClosed,
		}

		for i, priority := range ticketPriorities {
			for j, status := range ticketStatuses {
				// Create metadata as JSON
				metadata := map[string]interface{}{
					"source":     "Email",
					"resolution": "Fixed by restarting",
					"time_spent": "2 hours",
					"tags":       []string{"hardware", "urgent", "customer"},
				}

				metadataJSON, _ := json.Marshal(metadata)

				createdByID := staffMembers[i%len(staffMembers)].ID
				assignedToID := staffMembers[j%len(staffMembers)].ID
				categoryID := categories[i%len(categories)].ID
				assetID := assets[j%len(assets)].ID

				ticket := models.Ticket{
					TenantID:       tenant.ID,
					Title:          fmt.Sprintf("Ticket %d-%d: %s Priority Issue", i+1, j+1, priority),
					Description:    fmt.Sprintf("This is a sample ticket with %s priority and %s status", priority, status),
					Priority:       priority,
					Status:         status,
					CreatedByID:    &createdByID,
					AssignedToID:   &assignedToID,
					DueDate:        time.Now().AddDate(0, 0, i+j+1), // Due in i+j+1 days
					CategoryID:     &categoryID,
					RelatedAssetID: &assetID,
					Notes:          fmt.Sprintf("Sample notes for ticket %d-%d", i+1, j+1),
					Metadata:       string(metadataJSON),
				}

				result := database.DB.Create(&ticket)
				if result.Error != nil {
					log.Fatalf("Failed to create ticket: %v", result.Error)
				}
				fmt.Printf("Created ticket: %s (ID: %d) for tenant %s\n", ticket.Title, ticket.ID, tenant.Name)
			}
		}
	}

	fmt.Println("Seed data created successfully!")
}
