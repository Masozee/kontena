package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleSales   UserRole = "sales"
	RoleSupport UserRole = "support"
)

// User represents a user in the multi-tenant CRM system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TenantID  uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant    Tenant         `json:"-" gorm:"foreignKey:TenantID"` // Hide from JSON for Swagger
	Name      string         `json:"name" gorm:"size:100;not null"`
	Email     string         `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Role      UserRole       `json:"role" gorm:"size:20;not null"`
	Profile   string         `json:"profile" gorm:"type:jsonb"` // JSONB field for custom attributes
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Hide from JSON and Swagger
}
