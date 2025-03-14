package models

import (
	"time"

	"gorm.io/gorm"
)

// Person represents a team member in the project management system
type Person struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TenantID  uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant    *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Email     string         `json:"email" gorm:"size:100;not null;uniqueIndex:idx_tenant_email"`
	Role      string         `json:"role" gorm:"size:50;not null"`
	Position  string         `json:"position" gorm:"size:100"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Avatar    string         `json:"avatar" gorm:"size:500"`
	Projects  []*Project     `json:"projects,omitempty" gorm:"many2many:project_people;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
