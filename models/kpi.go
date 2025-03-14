package models

import (
	"time"

	"gorm.io/gorm"
)

// KPI represents a Key Performance Indicator for a project
type KPI struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProjectID    uint           `json:"project_id" gorm:"not null;index"`
	Project      *Project       `json:"-" gorm:"foreignKey:ProjectID"`
	Description  string         `json:"description" gorm:"type:text;not null"`
	TargetValue  float64        `json:"target_value" gorm:"not null"`
	CurrentValue float64        `json:"current_value" gorm:"not null;default:0"`
	Unit         string         `json:"unit" gorm:"size:20;not null"` // %, $, count, etc.
	Achieved     bool           `json:"achieved" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Progress calculates the KPI progress as a percentage
func (k *KPI) Progress() float64 {
	if k.TargetValue == 0 {
		return 0
	}
	return (k.CurrentValue / k.TargetValue) * 100
}

// UpdateAchievement sets Achieved to true if CurrentValue >= TargetValue
func (k *KPI) UpdateAchievement() {
	k.Achieved = k.CurrentValue >= k.TargetValue
}
