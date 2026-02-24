package core

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusQuote     OrderStatus = "QUOTE"
	OrderStatusSold      OrderStatus = "SOLD"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type MovementType string

const (
	MovementTypeAdjustment MovementType = "ADJUSTMENT"
	MovementTypeLoss       MovementType = "LOSS"
	MovementTypeDamaged    MovementType = "DAMAGED"
	MovementTypeReturn     MovementType = "RETURN"
	MovementTypeOther      MovementType = "OTHER"
)

type PurchaseStatus string

const (
	PurchaseStatusPending   PurchaseStatus = "PENDING"
	PurchaseStatusCompleted PurchaseStatus = "COMPLETED"
	PurchaseStatusCancelled PurchaseStatus = "CANCELLED"
)

type Product struct {
	ID           string    `json:"id" db:"id"`
	SKU          string    `json:"sku" db:"sku"`
	Name         string    `json:"name" db:"name"`
	Description  *string   `json:"description" db:"description"`
	CategoryID   *string   `json:"category_id" db:"category_id"`
	Price        float64   `json:"price" db:"price"`
	Cost         float64   `json:"cost" db:"cost"`
	CurrentStock int       `json:"current_stock" db:"current_stock"`
	MinStock     int       `json:"min_stock" db:"min_stock"`
	ImageURL     *string   `json:"image_url" db:"image_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	ID            string      `json:"id" db:"id"`
	CustomerID    *string     `json:"customer_id" db:"customer_id"`
	Status        OrderStatus `json:"status" db:"status"`
	TotalAmount   float64     `json:"total_amount" db:"total_amount"`
	Notes         *string     `json:"notes" db:"notes"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	ConvertedAt   *time.Time  `json:"converted_at" db:"converted_at"`
	ItemQuantity  int         `json:"item_quantity,omitempty" db:"item_quantity"`
	CustomerQuote *string     `json:"customer_quote,omitempty" db:"customer_quote"`
	InvoiceNumber *string     `json:"invoice_number" db:"invoice_number"`
	Items         []OrderItem `json:"items,omitempty" db:"-"`
}

type OrderItem struct {
	ID        string  `json:"id" db:"id"`
	OrderID   string  `json:"order_id" db:"order_id"`
	ProductID string  `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
}

type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Customer struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	TaxID     string    `json:"tax_id" db:"tax_id"` // RUT
	Email     *string   `json:"email" db:"email"`
	Phone     *string   `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Supplier struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	TaxID        string    `json:"tax_id" db:"tax_id"`
	ContactEmail *string   `json:"contact_email" db:"contact_email"`
	Phone        *string   `json:"phone" db:"phone"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Purchase struct {
	ID            string         `json:"id" db:"id"`
	SupplierID    *string        `json:"supplier_id" db:"supplier_id"`
	InvoiceNumber *string        `json:"invoice_number" db:"invoice_number"`
	InvoiceUrl    *string        `json:"invoice_url" db:"invoice_url"`
	TotalAmount   float64        `json:"total_amount" db:"total_amount"`
	Status        PurchaseStatus `json:"status" db:"status"`
	PurchaseDate  time.Time      `json:"purchase_date" db:"purchase_date"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	ItemQuantity  int            `json:"item_quantity,omitempty" db:"item_quantity"`
	// Relational field (There is not in the DB, it is filled in the JOIN or business logic)
	Items []PurchaseItem `json:"items,omitempty" db:"-"`
}

type PurchaseItem struct {
	ID         string    `json:"id" db:"id"`
	PurchaseID string    `json:"purchase_id" db:"purchase_id"`
	ProductID  string    `json:"product_id" db:"product_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	UnitCost   float64   `json:"unit_cost" db:"unit_cost"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type InventoryMovement struct {
	ID        string       `json:"id" db:"id"`
	ProductID string       `json:"product_id" db:"product_id"`
	Type      MovementType `json:"type" db:"type"`
	Quantity  int          `json:"quantity" db:"quantity"`
	Reason    string       `json:"reason" db:"reason"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}

//Auth

type User struct {
	ID           string `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}
