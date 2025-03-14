package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents a project in the project management system
type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant      *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Name        string         `json:"name" gorm:"size:200;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Budget      float64        `json:"budget"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Status      string         `json:"status" gorm:"size:20;not null;default:'planning'"`
	People      []*Person      `json:"people,omitempty" gorm:"many2many:project_people;"`
	KPIs        []KPI          `json:"kpis,omitempty" gorm:"foreignKey:ProjectID"`
	Tasks       []Task         `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
	Reports     []Report       `json:"reports,omitempty" gorm:"foreignKey:ProjectID"`
	Milestones  []Milestone    `json:"milestones,omitempty" gorm:"foreignKey:ProjectID"`
	Risks       []Risk         `json:"risks,omitempty" gorm:"foreignKey:ProjectID"`
	Issues      []Issue        `json:"issues,omitempty" gorm:"foreignKey:ProjectID"`
	Documents   []Document     `json:"documents,omitempty" gorm:"foreignKey:ProjectID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
