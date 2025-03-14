package models

import (
	"time"

	"gorm.io/gorm"
)

// Document represents a project document or attachment
type Document struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ProjectID    uint           `json:"project_id" gorm:"not null;index"`
	Project      *Project       `json:"-" gorm:"foreignKey:ProjectID"`
	Name         string         `json:"name" gorm:"size:200;not null"`
	FileURL      string         `json:"file_url" gorm:"size:500;not null"`
	UploadedByID uint           `json:"uploaded_by_id" gorm:"not null;index"`
	UploadedBy   *Person        `json:"uploaded_by" gorm:"foreignKey:UploadedByID"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
