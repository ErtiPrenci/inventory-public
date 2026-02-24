package core

import "time"

// Order
type CreateOrderRequest struct {
	CustomerID    *string              `json:"customer_id"`
	Status        OrderStatus          `json:"status"`
	Notes         *string              `json:"notes"`
	CustomerQuote *string              `json:"customer_quote,omitempty"`
	InvoiceNumber *string              `json:"invoice_number,omitempty"`
	Items         []CreateOrderItemDTO `json:"items"`
}

type CreateOrderItemDTO struct {
	ProductID string   `json:"product_id"`
	Quantity  int      `json:"quantity"`
	UnitPrice *float64 `json:"unit_price,omitempty"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	CustomerID    *string             `json:"customer_id"`
	Status        OrderStatus         `json:"status"`
	TotalAmount   float64             `json:"total_amount"`
	Notes         *string             `json:"notes"`
	CreatedAt     time.Time           `json:"created_at"`
	ConvertedAt   *time.Time          `json:"converted_at" db:"converted_at"`
	CustomerQuote *string             `json:"customer_quote" db:"customer_quote"`
	InvoiceNumber *string             `json:"invoice_number" db:"invoice_number"`
	Items         []OrderItemResponse `json:"items"`
}
type OrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	SubTotal    float64 `json:"subtotal"`
}

type UpdateOrderRequest struct {
	Status        OrderStatus `json:"status"`
	Notes         *string     `json:"notes,omitempty"`
	InvoiceNumber *string     `json:"invoice_number,omitempty"`
	CustomerID    *string     `json:"customer_id,omitempty"`
	CustomerQuote *string     `json:"customer_quote,omitempty"`
}

// Purchase
type CreatePurchaseRequest struct {
	SupplierID    *string              `json:"supplier_id"`
	InvoiceNumber *string              `json:"invoice_number"`
	InvoiceUrl    *string              `json:"invoice_url"`
	TotalAmount   float64              `json:"total_amount"`
	Status        PurchaseStatus       `json:"status"`
	PurchaseDate  time.Time            `json:"purchase_date"`
	NewItems      *bool                `json:"new_items"`
	Items         []CreatePurchaseItem `json:"items"`
}

type CreatePurchaseItem struct {
	ProductID *string  `json:"product_id"`
	Quantity  int      `json:"quantity"`
	UnitCost  float64  `json:"unit_cost"`
	Sku       *string  `json:"sku"`
	Name      *string  `json:"name"`
	Price     *float64 `json:"price"`
	MinStock  *int     `json:"min_stock"`
}

type UpdatePurchaseRequest struct {
	SupplierID    *string         `json:"supplier_id"`
	InvoiceNumber *string         `json:"invoice_number"`
	InvoiceUrl    *string         `json:"invoice_url"`
	Status        *PurchaseStatus `json:"status"`
	PurchaseDate  *time.Time      `json:"purchase_date"`
}

// Movements
type CreateMovementRequest struct {
	ProductID string       `json:"product_id"`
	Quantity  int          `json:"quantity"`
	Type      MovementType `json:"type"`
	Reason    string       `json:"reason"`
}

type InventoryMovementResponse struct {
	ID          string       `json:"id" db:"id"`
	ProductID   string       `json:"product_id" db:"product_id"`
	ProductName string       `json:"product_name" db:"product_name"`
	ProductSku  string       `json:"product_sku" db:"product_sku"`
	Quantity    int          `json:"quantity" db:"quantity"`
	Type        MovementType `json:"type" db:"type"`
	Reason      string       `json:"reason" db:"reason"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
}

type InventoryCreateMovementResponse struct {
	ID         string       `json:"id" db:"id"`
	ProductID  string       `json:"product_id" db:"product_id"`
	Type       MovementType `json:"type" db:"type"`
	Quantity   int          `json:"quantity" db:"quantity"`
	Reason     string       `json:"reason" db:"reason"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	FinalStock int          `json:"final_stock" db:"final_stock"`
}

// Pagination & Filters
type ProductFilter struct {
	Page       int
	Limit      int
	Search     string
	CategoryID string
}

type MetaData struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta MetaData    `json:"meta"`
}

type UploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	FileKey   string `json:"file_key"`
}
