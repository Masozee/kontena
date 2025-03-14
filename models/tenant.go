package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant represents a tenant in the multi-tenant project management system
type Tenant struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Plan        string         `json:"plan" gorm:"size:50;not null;default:'free'"`
	Status      string         `json:"status" gorm:"size:20;not null;default:'active'"`
	Domain      string         `json:"domain" gorm:"size:100;uniqueIndex"`
	LogoURL     string         `json:"logo_url" gorm:"size:500"`
	Settings    string         `json:"settings" gorm:"type:jsonb"` // JSONB field for tenant settings
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
