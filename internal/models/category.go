package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a category for leads with access control
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant      Tenant         `json:"-" gorm:"foreignKey:TenantID"` // Hide from JSON for Swagger
	Name        string         `json:"name" gorm:"size:100;not null"`
	Permissions string         `json:"permissions" gorm:"type:jsonb"` // JSONB field for permissions
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Hide from JSON and Swagger
}
