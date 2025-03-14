package models

import (
	"time"

	"gorm.io/gorm"
)

// StaffRole represents the role of a staff member in the system
type StaffRole string

const (
	RoleStaffAdmin    StaffRole = "admin"
	RoleStaffManager  StaffRole = "manager"
	RoleStaffEmployee StaffRole = "employee"
)

// Staff represents a staff member in the multi-tenant system
type Staff struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	TenantID   uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant     Tenant         `json:"-" gorm:"foreignKey:TenantID"` // Hide from JSON for Swagger
	Name       string         `json:"name" gorm:"size:100;not null"`
	Email      string         `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Role       StaffRole      `json:"role" gorm:"size:20;not null"`
	Department string         `json:"department" gorm:"size:50"`
	Position   string         `json:"position" gorm:"size:50"`
	Profile    string         `json:"profile" gorm:"type:jsonb"` // JSONB field for custom attributes
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"` // Hide from JSON and Swagger
}
