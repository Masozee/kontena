package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetStatus represents the status of an asset
type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "available"
	AssetStatusAssigned    AssetStatus = "assigned"
	AssetStatusMaintenance AssetStatus = "maintenance"
	AssetStatusRetired     AssetStatus = "retired"
)

// AssetType represents the type of an asset
type AssetType string

const (
	AssetTypeComputer  AssetType = "computer"
	AssetTypeVehicle   AssetType = "vehicle"
	AssetTypeFurniture AssetType = "furniture"
	AssetTypeEquipment AssetType = "equipment"
	AssetTypeOther     AssetType = "other"
)

// Asset represents a physical or digital asset in the system
type Asset struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	TenantID      uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant        Tenant         `json:"-" gorm:"foreignKey:TenantID"` // Hide from JSON for Swagger
	Name          string         `json:"name" gorm:"size:100;not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Type          AssetType      `json:"type" gorm:"size:20;not null"`
	SerialNumber  string         `json:"serial_number" gorm:"size:100"`
	PurchaseDate  time.Time      `json:"purchase_date"`
	PurchasePrice float64        `json:"purchase_price"`
	Status        AssetStatus    `json:"status" gorm:"size:20;not null"`
	AssignedToID  *uint          `json:"assigned_to_id" gorm:"index"`
	AssignedTo    Staff          `json:"-" gorm:"foreignKey:AssignedToID"` // Hide from JSON for Swagger
	Location      string         `json:"location" gorm:"size:100"`
	Notes         string         `json:"notes" gorm:"type:text"`
	Metadata      string         `json:"metadata" gorm:"type:jsonb"` // JSONB field for custom metadata
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"` // Hide from JSON and Swagger
}
