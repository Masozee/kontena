package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketPriority represents the priority of a ticket
type TicketPriority string

const (
	TicketPriorityLow      TicketPriority = "low"
	TicketPriorityMedium   TicketPriority = "medium"
	TicketPriorityHigh     TicketPriority = "high"
	TicketPriorityCritical TicketPriority = "critical"
)

// TicketStatus represents the status of a ticket
type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
)

// Ticket represents a support ticket in the system
type Ticket struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TenantID       uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant         Tenant         `json:"-" gorm:"foreignKey:TenantID"`
	Title          string         `json:"title" gorm:"size:200;not null"`
	Description    string         `json:"description" gorm:"type:text;not null"`
	Priority       TicketPriority `json:"priority" gorm:"size:20;not null"`
	Status         TicketStatus   `json:"status" gorm:"size:20;not null"`
	CreatedByID    *uint          `json:"created_by_id" gorm:"index"`
	CreatedBy      Staff          `json:"-" gorm:"foreignKey:CreatedByID"`
	AssignedToID   *uint          `json:"assigned_to_id" gorm:"index"`
	AssignedTo     Staff          `json:"-" gorm:"foreignKey:AssignedToID"`
	DueDate        time.Time      `json:"due_date"`
	CategoryID     *uint          `json:"category_id" gorm:"index"`
	Category       Category       `json:"-" gorm:"foreignKey:CategoryID"`
	RelatedAssetID *uint          `json:"related_asset_id" gorm:"index"`
	RelatedAsset   Asset          `json:"-" gorm:"foreignKey:RelatedAssetID"`
	Notes          string         `json:"notes" gorm:"type:text"`
	Metadata       string         `json:"metadata" gorm:"type:jsonb"` // JSONB field for custom metadata
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"` // Hide from JSON and Swagger
}
