package models

import (
	"time"

	"gorm.io/gorm"
)

// RiskStatus represents the status of a risk
type RiskStatus string

const (
	RiskStatusIdentified RiskStatus = "identified"
	RiskStatusMonitoring RiskStatus = "monitoring"
	RiskStatusMitigated  RiskStatus = "mitigated"
	RiskStatusClosed     RiskStatus = "closed"
)

// Risk represents a project risk
type Risk struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProjectID   uint           `json:"project_id" gorm:"not null;index"`
	Project     *Project       `json:"-" gorm:"foreignKey:ProjectID"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Impact      string         `json:"impact" gorm:"size:50;not null"`      // High, Medium, Low
	Probability string         `json:"probability" gorm:"size:50;not null"` // High, Medium, Low
	Mitigation  string         `json:"mitigation" gorm:"type:text"`
	Status      RiskStatus     `json:"status" gorm:"size:20;not null;default:'identified'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
