package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetCategory represents a category of assets
type AssetCategory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant      *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	ParentID    *uint          `json:"parent_id" gorm:"index"`
	Parent      *AssetCategory `json:"-" gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// AssetStatus represents the status of an asset
type AssetStatus string

const (
	AssetStatusProcurement AssetStatus = "procurement"
	AssetStatusInStock     AssetStatus = "in_stock"
	AssetStatusAssigned    AssetStatus = "assigned"
	AssetStatusMaintenance AssetStatus = "maintenance"
	AssetStatusRetired     AssetStatus = "retired"
)

// Asset represents a physical or digital asset
type Asset struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	TenantID         uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant           *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Name             string         `json:"name" gorm:"size:100;not null"`
	Description      string         `json:"description" gorm:"type:text"`
	CategoryID       uint           `json:"category_id" gorm:"not null;index"`
	Category         *AssetCategory `json:"category" gorm:"foreignKey:CategoryID"`
	SerialNumber     string         `json:"serial_number" gorm:"size:100"`
	ModelNumber      string         `json:"model_number" gorm:"size:100"`
	Manufacturer     string         `json:"manufacturer" gorm:"size:100"`
	PurchaseDate     *time.Time     `json:"purchase_date"`
	PurchasePrice    float64        `json:"purchase_price"`
	WarrantyExpiry   *time.Time     `json:"warranty_expiry"`
	Status           AssetStatus    `json:"status" gorm:"size:20;not null;default:'in_stock'"`
	LocationID       *uint          `json:"location_id" gorm:"index"`
	Location         *Location      `json:"location" gorm:"foreignKey:LocationID"`
	CurrentAssignee  *uint          `json:"current_assignee" gorm:"index"`
	AssignedTo       *Person        `json:"assigned_to" gorm:"foreignKey:CurrentAssignee"`
	ExpectedLifespan int            `json:"expected_lifespan"` // in months
	Notes            string         `json:"notes" gorm:"type:text"`
	Tags             string         `json:"tags" gorm:"size:255"`
	Barcode          string         `json:"barcode" gorm:"size:100"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// Location represents a physical location where assets can be stored
type Location struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TenantID    uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant      *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Address     string         `json:"address" gorm:"type:text"`
	Type        string         `json:"type" gorm:"size:50"` // warehouse, office, etc.
	ParentID    *uint          `json:"parent_id" gorm:"index"`
	Parent      *Location      `json:"-" gorm:"foreignKey:ParentID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Vendor represents a supplier or service provider
type Vendor struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	TenantID     uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant       *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Name         string         `json:"name" gorm:"size:100;not null"`
	ContactName  string         `json:"contact_name" gorm:"size:100"`
	ContactEmail string         `json:"contact_email" gorm:"size:100"`
	ContactPhone string         `json:"contact_phone" gorm:"size:20"`
	Address      string         `json:"address" gorm:"type:text"`
	Website      string         `json:"website" gorm:"size:255"`
	Notes        string         `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// ProcurementStatus represents the status of a procurement request
type ProcurementStatus string

const (
	ProcurementStatusDraft     ProcurementStatus = "draft"
	ProcurementStatusSubmitted ProcurementStatus = "submitted"
	ProcurementStatusApproved  ProcurementStatus = "approved"
	ProcurementStatusRejected  ProcurementStatus = "rejected"
	ProcurementStatusOrdered   ProcurementStatus = "ordered"
	ProcurementStatusReceived  ProcurementStatus = "received"
	ProcurementStatusCancelled ProcurementStatus = "cancelled"
)

// ProcurementRequest represents a request to procure assets
type ProcurementRequest struct {
	ID             uint              `json:"id" gorm:"primaryKey"`
	TenantID       uint              `json:"tenant_id" gorm:"not null;index"`
	Tenant         *Tenant           `json:"-" gorm:"foreignKey:TenantID"`
	RequestNumber  string            `json:"request_number" gorm:"size:50;not null"`
	RequestedByID  uint              `json:"requested_by_id" gorm:"not null;index"`
	RequestedBy    *Person           `json:"requested_by" gorm:"foreignKey:RequestedByID"`
	ApprovedByID   *uint             `json:"approved_by_id" gorm:"index"`
	ApprovedBy     *Person           `json:"approved_by" gorm:"foreignKey:ApprovedByID"`
	Status         ProcurementStatus `json:"status" gorm:"size:20;not null;default:'draft'"`
	RequestDate    time.Time         `json:"request_date" gorm:"not null"`
	ApprovalDate   *time.Time        `json:"approval_date"`
	ExpectedDate   *time.Time        `json:"expected_date"`
	TotalBudget    float64           `json:"total_budget"`
	Notes          string            `json:"notes" gorm:"type:text"`
	Items          []ProcurementItem `json:"items" gorm:"foreignKey:ProcurementID"`
	PurchaseOrders []PurchaseOrder   `json:"purchase_orders" gorm:"foreignKey:ProcurementID"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `json:"-" gorm:"index"`
}

// ProcurementItem represents an item in a procurement request
type ProcurementItem struct {
	ID                uint                `json:"id" gorm:"primaryKey"`
	ProcurementID     uint                `json:"procurement_id" gorm:"not null;index"`
	Procurement       *ProcurementRequest `json:"-" gorm:"foreignKey:ProcurementID"`
	CategoryID        uint                `json:"category_id" gorm:"not null;index"`
	Category          *AssetCategory      `json:"category" gorm:"foreignKey:CategoryID"`
	Description       string              `json:"description" gorm:"type:text;not null"`
	Quantity          int                 `json:"quantity" gorm:"not null"`
	EstimatedPrice    float64             `json:"estimated_price"`
	PreferredVendorID *uint               `json:"preferred_vendor_id" gorm:"index"`
	PreferredVendor   *Vendor             `json:"preferred_vendor" gorm:"foreignKey:PreferredVendorID"`
	Justification     string              `json:"justification" gorm:"type:text"`
	Status            string              `json:"status" gorm:"size:20;not null;default:'pending'"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	DeletedAt         gorm.DeletedAt      `json:"-" gorm:"index"`
}

// PurchaseOrderStatus represents the status of a purchase order
type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft     PurchaseOrderStatus = "draft"
	PurchaseOrderStatusSent      PurchaseOrderStatus = "sent"
	PurchaseOrderStatusConfirmed PurchaseOrderStatus = "confirmed"
	PurchaseOrderStatusPartial   PurchaseOrderStatus = "partially_received"
	PurchaseOrderStatusReceived  PurchaseOrderStatus = "received"
	PurchaseOrderStatusCancelled PurchaseOrderStatus = "cancelled"
)

// PurchaseOrder represents an order to purchase assets
type PurchaseOrder struct {
	ID               uint                `json:"id" gorm:"primaryKey"`
	TenantID         uint                `json:"tenant_id" gorm:"not null;index"`
	Tenant           *Tenant             `json:"-" gorm:"foreignKey:TenantID"`
	OrderNumber      string              `json:"order_number" gorm:"size:50;not null"`
	ProcurementID    *uint               `json:"procurement_id" gorm:"index"`
	Procurement      *ProcurementRequest `json:"-" gorm:"foreignKey:ProcurementID"`
	VendorID         uint                `json:"vendor_id" gorm:"not null;index"`
	Vendor           *Vendor             `json:"vendor" gorm:"foreignKey:VendorID"`
	Status           PurchaseOrderStatus `json:"status" gorm:"size:20;not null;default:'draft'"`
	OrderDate        time.Time           `json:"order_date" gorm:"not null"`
	ExpectedDelivery *time.Time          `json:"expected_delivery"`
	DeliveryAddress  string              `json:"delivery_address" gorm:"type:text"`
	ShippingCost     float64             `json:"shipping_cost"`
	TaxAmount        float64             `json:"tax_amount"`
	TotalAmount      float64             `json:"total_amount"`
	PaymentTerms     string              `json:"payment_terms" gorm:"size:255"`
	Notes            string              `json:"notes" gorm:"type:text"`
	Items            []PurchaseOrderItem `json:"items" gorm:"foreignKey:PurchaseOrderID"`
	Receipts         []AssetReceipt      `json:"receipts" gorm:"foreignKey:PurchaseOrderID"`
	CreatedByID      uint                `json:"created_by_id" gorm:"not null;index"`
	CreatedBy        *Person             `json:"created_by" gorm:"foreignKey:CreatedByID"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	DeletedAt        gorm.DeletedAt      `json:"-" gorm:"index"`
}

// PurchaseOrderItem represents an item in a purchase order
type PurchaseOrderItem struct {
	ID                uint             `json:"id" gorm:"primaryKey"`
	PurchaseOrderID   uint             `json:"purchase_order_id" gorm:"not null;index"`
	PurchaseOrder     *PurchaseOrder   `json:"-" gorm:"foreignKey:PurchaseOrderID"`
	ProcurementItemID *uint            `json:"procurement_item_id" gorm:"index"`
	ProcurementItem   *ProcurementItem `json:"-" gorm:"foreignKey:ProcurementItemID"`
	Description       string           `json:"description" gorm:"type:text;not null"`
	Quantity          int              `json:"quantity" gorm:"not null"`
	UnitPrice         float64          `json:"unit_price" gorm:"not null"`
	TotalPrice        float64          `json:"total_price" gorm:"not null"`
	ReceivedQuantity  int              `json:"received_quantity" gorm:"default:0"`
	ExpectedDate      *time.Time       `json:"expected_date"`
	Notes             string           `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `json:"-" gorm:"index"`
}

// AssetReceipt represents the receipt of assets from a purchase order
type AssetReceipt struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	TenantID        uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant          *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	ReceiptNumber   string         `json:"receipt_number" gorm:"size:50;not null"`
	PurchaseOrderID uint           `json:"purchase_order_id" gorm:"not null;index"`
	PurchaseOrder   *PurchaseOrder `json:"-" gorm:"foreignKey:PurchaseOrderID"`
	ReceivedByID    uint           `json:"received_by_id" gorm:"not null;index"`
	ReceivedBy      *Person        `json:"received_by" gorm:"foreignKey:ReceivedByID"`
	ReceiptDate     time.Time      `json:"receipt_date" gorm:"not null"`
	DeliveryNote    string         `json:"delivery_note" gorm:"size:100"`
	Notes           string         `json:"notes" gorm:"type:text"`
	Items           []ReceiptItem  `json:"items" gorm:"foreignKey:ReceiptID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// ReceiptItem represents an item in an asset receipt
type ReceiptItem struct {
	ID                  uint               `json:"id" gorm:"primaryKey"`
	ReceiptID           uint               `json:"receipt_id" gorm:"not null;index"`
	Receipt             *AssetReceipt      `json:"-" gorm:"foreignKey:ReceiptID"`
	PurchaseOrderItemID uint               `json:"purchase_order_item_id" gorm:"not null;index"`
	PurchaseOrderItem   *PurchaseOrderItem `json:"-" gorm:"foreignKey:PurchaseOrderItemID"`
	QuantityReceived    int                `json:"quantity_received" gorm:"not null"`
	Notes               string             `json:"notes" gorm:"type:text"`
	Assets              []Asset            `json:"assets" gorm:"many2many:receipt_item_assets;"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	DeletedAt           gorm.DeletedAt     `json:"-" gorm:"index"`
}

// MaintenanceType represents the type of maintenance
type MaintenanceType string

const (
	MaintenanceTypePreventive  MaintenanceType = "preventive"
	MaintenanceTypeCorrective  MaintenanceType = "corrective"
	MaintenanceTypeCalibration MaintenanceType = "calibration"
	MaintenanceTypeInspection  MaintenanceType = "inspection"
)

// MaintenanceStatus represents the status of a maintenance record
type MaintenanceStatus string

const (
	MaintenanceStatusScheduled  MaintenanceStatus = "scheduled"
	MaintenanceStatusInProgress MaintenanceStatus = "in_progress"
	MaintenanceStatusCompleted  MaintenanceStatus = "completed"
	MaintenanceStatusCancelled  MaintenanceStatus = "cancelled"
)

// MaintenanceRecord represents a maintenance activity for an asset
type MaintenanceRecord struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	TenantID        uint              `json:"tenant_id" gorm:"not null;index"`
	Tenant          *Tenant           `json:"-" gorm:"foreignKey:TenantID"`
	AssetID         uint              `json:"asset_id" gorm:"not null;index"`
	Asset           *Asset            `json:"asset" gorm:"foreignKey:AssetID"`
	MaintenanceType MaintenanceType   `json:"maintenance_type" gorm:"size:20;not null"`
	Status          MaintenanceStatus `json:"status" gorm:"size:20;not null;default:'scheduled'"`
	ScheduledDate   time.Time         `json:"scheduled_date" gorm:"not null"`
	CompletedDate   *time.Time        `json:"completed_date"`
	PerformedByID   *uint             `json:"performed_by_id" gorm:"index"`
	PerformedBy     *Person           `json:"performed_by" gorm:"foreignKey:PerformedByID"`
	VendorID        *uint             `json:"vendor_id" gorm:"index"`
	Vendor          *Vendor           `json:"vendor" gorm:"foreignKey:VendorID"`
	Cost            float64           `json:"cost"`
	Description     string            `json:"description" gorm:"type:text;not null"`
	Results         string            `json:"results" gorm:"type:text"`
	NextScheduled   *time.Time        `json:"next_scheduled"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"-" gorm:"index"`
}

// AssetAssignment represents the assignment of an asset to a person
type AssetAssignment struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TenantID       uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant         *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	AssetID        uint           `json:"asset_id" gorm:"not null;index"`
	Asset          *Asset         `json:"asset" gorm:"foreignKey:AssetID"`
	AssignedToID   uint           `json:"assigned_to_id" gorm:"not null;index"`
	AssignedTo     *Person        `json:"assigned_to" gorm:"foreignKey:AssignedToID"`
	AssignedByID   uint           `json:"assigned_by_id" gorm:"not null;index"`
	AssignedBy     *Person        `json:"assigned_by" gorm:"foreignKey:AssignedByID"`
	AssignmentDate time.Time      `json:"assignment_date" gorm:"not null"`
	ReturnDate     *time.Time     `json:"return_date"`
	ExpectedReturn *time.Time     `json:"expected_return"`
	Notes          string         `json:"notes" gorm:"type:text"`
	Status         string         `json:"status" gorm:"size:20;not null;default:'active'"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// StockTransaction represents a movement of assets in or out of inventory
type StockTransaction struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	TenantID         uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant           *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	TransactionType  string         `json:"transaction_type" gorm:"size:20;not null"` // in, out, transfer
	SourceLocationID *uint          `json:"source_location_id" gorm:"index"`
	SourceLocation   *Location      `json:"source_location" gorm:"foreignKey:SourceLocationID"`
	DestLocationID   *uint          `json:"dest_location_id" gorm:"index"`
	DestLocation     *Location      `json:"dest_location" gorm:"foreignKey:DestLocationID"`
	TransactionDate  time.Time      `json:"transaction_date" gorm:"not null"`
	PerformedByID    uint           `json:"performed_by_id" gorm:"not null;index"`
	PerformedBy      *Person        `json:"performed_by" gorm:"foreignKey:PerformedByID"`
	ReferenceNumber  string         `json:"reference_number" gorm:"size:50"`
	Notes            string         `json:"notes" gorm:"type:text"`
	Items            []StockItem    `json:"items" gorm:"foreignKey:TransactionID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// StockItem represents an item in a stock transaction
type StockItem struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	TransactionID uint              `json:"transaction_id" gorm:"not null;index"`
	Transaction   *StockTransaction `json:"-" gorm:"foreignKey:TransactionID"`
	AssetID       uint              `json:"asset_id" gorm:"not null;index"`
	Asset         *Asset            `json:"asset" gorm:"foreignKey:AssetID"`
	Quantity      int               `json:"quantity" gorm:"not null"`
	Notes         string            `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
}

// InventoryCount represents a physical inventory count
type InventoryCount struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	TenantID      uint           `json:"tenant_id" gorm:"not null;index"`
	Tenant        *Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	LocationID    uint           `json:"location_id" gorm:"not null;index"`
	Location      *Location      `json:"location" gorm:"foreignKey:LocationID"`
	CountDate     time.Time      `json:"count_date" gorm:"not null"`
	PerformedByID uint           `json:"performed_by_id" gorm:"not null;index"`
	PerformedBy   *Person        `json:"performed_by" gorm:"foreignKey:PerformedByID"`
	Status        string         `json:"status" gorm:"size:20;not null;default:'draft'"` // draft, in_progress, completed
	Notes         string         `json:"notes" gorm:"type:text"`
	Items         []CountItem    `json:"items" gorm:"foreignKey:InventoryCountID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// CountItem represents an item in an inventory count
type CountItem struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	InventoryCountID uint            `json:"inventory_count_id" gorm:"not null;index"`
	InventoryCount   *InventoryCount `json:"-" gorm:"foreignKey:InventoryCountID"`
	AssetID          uint            `json:"asset_id" gorm:"not null;index"`
	Asset            *Asset          `json:"asset" gorm:"foreignKey:AssetID"`
	ExpectedQuantity int             `json:"expected_quantity"`
	ActualQuantity   int             `json:"actual_quantity"`
	Discrepancy      int             `json:"discrepancy"`
	Notes            string          `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `json:"-" gorm:"index"`
}
