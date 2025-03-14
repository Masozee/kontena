package models

import (
	"time"

	"gorm.io/gorm"
)

// MilestoneStatus represents the status of a milestone
type MilestoneStatus string

const (
	MilestoneStatusPlanned    MilestoneStatus = "planned"
	MilestoneStatusInProgress MilestoneStatus = "in_progress"
	MilestoneStatusCompleted  MilestoneStatus = "completed"
	MilestoneStatusDelayed    MilestoneStatus = "delayed"
)

// Milestone represents a project milestone
type Milestone struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	ProjectID   uint            `json:"project_id" gorm:"not null;index"`
	Project     *Project        `json:"-" gorm:"foreignKey:ProjectID"`
	Title       string          `json:"title" gorm:"size:200;not null"`
	Description string          `json:"description" gorm:"type:text"`
	DueDate     time.Time       `json:"due_date" gorm:"not null"`
	Status      MilestoneStatus `json:"status" gorm:"size:20;not null;default:'planned'"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"index"`
}
