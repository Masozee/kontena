package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant represents a tenant in the multi-tenant CRM system
type Tenant struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Plan      string         `json:"plan" gorm:"size:50;not null"`
	Status    string         `json:"status" gorm:"size:20;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
