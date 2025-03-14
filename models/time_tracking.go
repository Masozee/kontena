package models

import (
	"time"

	"gorm.io/gorm"
)

// TimeTracking represents time spent on a task
type TimeTracking struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TaskID    uint           `json:"task_id" gorm:"not null;index"`
	Task      *Task          `json:"-" gorm:"foreignKey:TaskID"`
	PersonID  uint           `json:"person_id" gorm:"not null;index"`
	Person    *Person        `json:"person" gorm:"foreignKey:PersonID"`
	Hours     float64        `json:"hours" gorm:"not null"`
	Date      time.Time      `json:"date" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
