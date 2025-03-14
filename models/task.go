package models

import (
	"time"

	"gorm.io/gorm"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusBlocked    TaskStatus = "blocked"
)

// Task represents a task in a project
type Task struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProjectID    uint           `json:"project_id" gorm:"not null;index"`
	Project      *Project       `json:"-" gorm:"foreignKey:ProjectID"`
	Title        string         `json:"title" gorm:"size:200;not null"`
	Description  string         `json:"description" gorm:"type:text"`
	AssignedToID *uint          `json:"assigned_to_id" gorm:"index"`
	AssignedTo   *Person        `json:"assigned_to" gorm:"foreignKey:AssignedToID"`
	Status       TaskStatus     `json:"status" gorm:"size:20;not null;default:'todo'"`
	DueDate      *time.Time     `json:"due_date"`
	TimeEntries  []TimeTracking `json:"time_entries,omitempty" gorm:"foreignKey:TaskID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
