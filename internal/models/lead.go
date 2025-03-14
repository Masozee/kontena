package models

import (
	"time"

	"gorm.io/gorm"
)

// Lead represents a lead in the CRM system
type Lead struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	TenantID     uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant       Tenant         `json:"-" gorm:"foreignKey:TenantID"`
	Name         string         `json:"name" gorm:"size:100;not null"`
	Contact      string         `json:"contact" gorm:"size:100"`
	Email        string         `json:"email" gorm:"size:100"`
	CategoryID   *uint          `json:"category_id" gorm:"index"`
	Category     Category       `json:"-" gorm:"foreignKey:CategoryID"`
	AssignedTo   *uint          `json:"assigned_to" gorm:"index"`
	AssignedUser User           `json:"-" gorm:"foreignKey:AssignedTo"`
	Status       string         `json:"status" gorm:"size:50"`
	Notes        string         `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
