package models

import (
	"time"

	"gorm.io/gorm"
)

// IssueStatus represents the status of an issue
type IssueStatus string

const (
	IssueStatusOpen       IssueStatus = "open"
	IssueStatusInProgress IssueStatus = "in_progress"
	IssueStatusResolved   IssueStatus = "resolved"
	IssueStatusClosed     IssueStatus = "closed"
)

// Issue represents a project issue or ticket
type Issue struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProjectID    uint           `json:"project_id" gorm:"not null;index"`
	Project      *Project       `json:"-" gorm:"foreignKey:ProjectID"`
	Title        string         `json:"title" gorm:"size:200;not null"`
	Description  string         `json:"description" gorm:"type:text;not null"`
	Status       IssueStatus    `json:"status" gorm:"size:20;not null;default:'open'"`
	ReportedByID uint           `json:"reported_by_id" gorm:"not null;index"`
	ReportedBy   *Person        `json:"reported_by" gorm:"foreignKey:ReportedByID"`
	AssignedToID *uint          `json:"assigned_to_id" gorm:"index"`
	AssignedTo   *Person        `json:"assigned_to" gorm:"foreignKey:AssignedToID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
