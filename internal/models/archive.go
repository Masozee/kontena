package models

import (
	"time"

	"gorm.io/gorm"
)

// ArchiveStatus represents the status of an archive
type ArchiveStatus string

const (
	ArchiveStatusActive       ArchiveStatus = "active"
	ArchiveStatusInactive     ArchiveStatus = "inactive"
	ArchiveStatusConfidential ArchiveStatus = "confidential"
)

// Archive represents a document or file archive in the system
type Archive struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant      Tenant         `json:"-" gorm:"foreignKey:TenantID"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Description string         `json:"description" gorm:"type:text"`
	CategoryID  *uint          `json:"category_id" gorm:"index"`
	Category    Category       `json:"-" gorm:"foreignKey:CategoryID"`
	FilePath    string         `json:"file_path" gorm:"size:500"`
	FileType    string         `json:"file_type" gorm:"size:50"`
	FileSize    int64          `json:"file_size"`
	Status      ArchiveStatus  `json:"status" gorm:"size:20;not null"`
	CreatedByID *uint          `json:"created_by_id" gorm:"index"`
	CreatedBy   Staff          `json:"-" gorm:"foreignKey:CreatedByID"`
	Metadata    string         `json:"metadata" gorm:"type:jsonb"` // JSONB field for custom metadata
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Hide from JSON and Swagger
}
