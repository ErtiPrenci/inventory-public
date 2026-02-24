package repository

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-backend/internal/core"
	"log"
	"time"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]core.Order, error)
	GetByID(ctx context.Context, id string) (core.OrderResponse, error)
	Create(ctx context.Context, o core.Order) (core.Order, error)
	Update(ctx context.Context, o core.Order) (core.Order, error)
}

type postgresOrderRepo struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepo{db: db}
}

func (p *postgresOrderRepo) GetAll(ctx context.Context) ([]core.Order, error) {
	query := `
        SELECT 
			o.id, 
			o.customer_id, 
			o.status, 
			o.total_amount,
			o.notes,
			o.created_at,
			o.converted_at,
			o.customer_quote,
			o.invoice_number,
			COUNT(i.id) AS item_quantity
		FROM orders o
		LEFT JOIN order_items i ON i.order_id = o.id
		GROUP BY o.id
		ORDER BY o.created_at DESC;
    `

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]core.Order, 0)

	for rows.Next() {
		var o core.Order

		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.Status,
			&o.TotalAmount,
			&o.Notes,
			&o.CreatedAt,
			&o.ConvertedAt,
			&o.CustomerQuote,
			&o.InvoiceNumber,
			&o.ItemQuantity,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (p *postgresOrderRepo) GetByID(ctx context.Context, id string) (core.OrderResponse, error) {
	var order core.OrderResponse

	// Get the Header (Order)
	queryOrder := `
        SELECT 
            id, 
            customer_id, 
            total_amount, 
            status, 
            notes, 
            created_at,
			converted_at,
			customer_quote,
			invoice_number
        FROM orders
        WHERE id = $1
		LIMIT 30
		`

	// Scan the basic data
	err := p.db.QueryRowContext(ctx, queryOrder, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.TotalAmount,
		&order.Status,
		&order.Notes,
		&order.CreatedAt,
		&order.ConvertedAt,
		&order.CustomerQuote,
		&order.InvoiceNumber,
	)
	if err != nil {
		log.Printf("Failed in Header, Error %v", err)
		return core.OrderResponse{}, err
	}

	// Get the Details (Items + Product Data)
	queryItems := `
        SELECT 
            oi.id, 
            oi.product_id, 
            p.name,
            p.sku,
            oi.quantity, 
            oi.unit_price
        FROM order_items oi
        JOIN products p ON oi.product_id = p.id
        WHERE oi.order_id = $1`

	rows, err := p.db.QueryContext(ctx, queryItems, id)
	if err != nil {
		log.Printf("Failed in Details, Error %v", err)
		return core.OrderResponse{}, err
	}
	defer rows.Close()

	// Initialize the array empty so that JSON shows [] and not null if there are no items
	order.Items = make([]core.OrderItemResponse, 0)

	for rows.Next() {
		var item core.OrderItemResponse

		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.SKU,
			&item.Quantity,
			&item.UnitPrice,
		)
		if err != nil {
			log.Printf("Failed in filling fields, Error %v", err)
			return core.OrderResponse{}, err
		}

		// Calculate the subtotal on the fly for frontend convenience
		item.SubTotal = item.UnitPrice * float64(item.Quantity)

		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (p *postgresOrderRepo) Create(ctx context.Context, o core.Order) (core.Order, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin tx: %v", err)
		return core.Order{}, err
	}
	defer tx.Rollback()

	// Insert Header
	queryOrder := `
        INSERT INTO orders (customer_id, status, total_amount, notes, customer_quote, invoice_number)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, converted_at
    `
	err = tx.QueryRowContext(ctx, queryOrder, o.CustomerID, o.Status, o.TotalAmount, o.Notes, o.CustomerQuote, o.InvoiceNumber).Scan(&o.ID, &o.CreatedAt, &o.ConvertedAt)
	if err != nil {
		log.Printf("insert order: %v", err)
		return core.Order{}, fmt.Errorf("insert order: %w", err)
	}

	// Prepare Statements
	queryItem := `INSERT INTO order_items (order_id, product_id, quantity, unit_price, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	stmtItem, err := tx.PrepareContext(ctx, queryItem)
	if err != nil {
		log.Printf("prepare item: %v", err)
		return core.Order{}, fmt.Errorf("prepare item: %w", err)
	}
	defer stmtItem.Close()

	// Query to deduct stock (Only used if it's SOLD)
	queryStock := `UPDATE products SET current_stock = current_stock - $1 WHERE id = $2 RETURNING current_stock`
	stmtStock, err := tx.PrepareContext(ctx, queryStock)
	if err != nil {
		log.Printf("prepare stock: %v", err)
		return core.Order{}, fmt.Errorf("prepare stock: %w", err)
	}
	defer stmtStock.Close()

	for i, item := range o.Items {
		// Insert Item
		var itemID string
		err := stmtItem.QueryRowContext(ctx, o.ID, item.ProductID, item.Quantity, item.UnitPrice).Scan(&itemID)
		if err != nil {
			log.Printf("insert item: %v", err)
			return core.Order{}, fmt.Errorf("insert item: %w", err)
		}
		o.Items[i].ID = itemID

		// Deduct stock (Only if it's SOLD)
		if o.Status == core.OrderStatusSold {
			var newStock int
			err := stmtStock.QueryRowContext(ctx, item.Quantity, item.ProductID).Scan(&newStock)
			if err != nil {
				log.Printf("error deducting stock for product %s: %v", item.ProductID, err)
				return core.Order{}, fmt.Errorf("error deducting stock for product %s: %w", item.ProductID, err)
			}
			if newStock < 0 {
				log.Printf("insufficient stock for product %s", item.ProductID)
				return core.Order{}, fmt.Errorf("insufficient stock for product ID %s", item.ProductID)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error committing transaction: %v", err)
		return core.Order{}, err
	}
	return o, nil
}

func (p *postgresOrderRepo) Update(ctx context.Context, o core.Order) (core.Order, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return core.Order{}, err
	}
	defer tx.Rollback()

	var dbStatus core.OrderStatus
	var dbConvertedAt *time.Time
	var dbTotalAmount float64
	var dbNotes, dbCustomerID, dbInvoiceNumber, dbCustomerQuote *string
	var dbCreatedAt time.Time

	queryFetch := `SELECT status, converted_at, total_amount, notes, customer_id, invoice_number, customer_quote, created_at FROM orders WHERE id = $1 FOR UPDATE`

	err = tx.QueryRowContext(ctx, queryFetch, o.ID).Scan(
		&dbStatus,
		&dbConvertedAt,
		&dbTotalAmount,
		&dbNotes,
		&dbCustomerID,
		&dbInvoiceNumber,
		&dbCustomerQuote,
		&dbCreatedAt,
	)
	if err != nil {
		return core.Order{}, fmt.Errorf("error fetching real order: %w", err)
	}

	var finalConvertedAt *time.Time
	finalConvertedAt = dbConvertedAt

	// When turn to SOLD, overwrite with NOW()
	if dbStatus != core.OrderStatusSold && o.Status == core.OrderStatusSold {
		now := time.Now()
		finalConvertedAt = &now
	}

	// UPDATE HEADER - DYNAMIC
	queryUpdate := "UPDATE orders SET "
	args := []interface{}{}
	argIdx := 1
	if o.Status != "" {
		queryUpdate += fmt.Sprintf("status = $%d, ", argIdx)
		args = append(args, o.Status)
		argIdx++
	}
	if o.Notes != nil {
		if *o.Notes == "" {
			o.Notes = nil
		}
		queryUpdate += fmt.Sprintf("notes = $%d, ", argIdx)
		args = append(args, o.Notes)
		argIdx++
	}
	if o.CustomerID != nil {
		if *o.CustomerID == "" {
			o.CustomerID = nil
		}
		queryUpdate += fmt.Sprintf("customer_id = $%d, ", argIdx)
		args = append(args, o.CustomerID)
		argIdx++
	}
	if o.InvoiceNumber != nil {
		if *o.InvoiceNumber == "" {
			o.InvoiceNumber = nil
		}
		queryUpdate += fmt.Sprintf("invoice_number = $%d, ", argIdx)
		args = append(args, o.InvoiceNumber)
		argIdx++
	}
	if o.CustomerQuote != nil {
		if *o.CustomerQuote == "" {
			o.CustomerQuote = nil
		}
		queryUpdate += fmt.Sprintf("customer_quote = $%d, ", argIdx)
		args = append(args, o.CustomerQuote)
		argIdx++
	}

	queryUpdate += fmt.Sprintf("converted_at = $%d ", argIdx)
	args = append(args, finalConvertedAt)
	argIdx++

	queryUpdate += fmt.Sprintf("WHERE id = $%d", argIdx)
	args = append(args, o.ID)

	_, err = tx.ExecContext(ctx, queryUpdate, args...)
	if err != nil {
		return core.Order{}, fmt.Errorf("update header: %w", err)
	}

	// STOCK LOGIC
	if o.Status != "" && dbStatus != o.Status {
		//TODO:[ ] Add Validation for SOLD to QUOTE
		if dbStatus == core.OrderStatusQuote && o.Status == core.OrderStatusSold {
			if err := p.adjustStock(ctx, tx, o.ID, -1); err != nil {
				return core.Order{}, err
			}
		}
		if dbStatus == core.OrderStatusSold && o.Status == core.OrderStatusCancelled {
			if err := p.adjustStock(ctx, tx, o.ID, 1); err != nil {
				return core.Order{}, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return core.Order{}, err
	}

	// Update return struct for frontend with FULL object
	o.TotalAmount = dbTotalAmount
	o.CreatedAt = dbCreatedAt
	o.ConvertedAt = finalConvertedAt

	if o.Status == "" {
		o.Status = dbStatus
	}
	if o.Notes == nil {
		o.Notes = dbNotes
	}
	if o.CustomerID == nil {
		o.CustomerID = dbCustomerID
	}
	if o.InvoiceNumber == nil {
		o.InvoiceNumber = dbInvoiceNumber
	}
	if o.CustomerQuote == nil {
		o.CustomerQuote = dbCustomerQuote
	}

	return o, nil
}

func (p *postgresOrderRepo) adjustStock(ctx context.Context, tx *sql.Tx, orderID string, multiplier int) error {
	// First deduct/add the stock
	query := `
        UPDATE products p
        SET current_stock = p.current_stock + (oi.quantity * $1)
        FROM order_items oi
        WHERE oi.product_id = p.id AND oi.order_id = $2
		RETURNING p.id, p.current_stock
    `
	rows, err := tx.QueryContext(ctx, query, multiplier, orderID)
	if err != nil {
		return fmt.Errorf("adjust stock error: %w", err)
	}
	defer rows.Close()

	// Check if any product went below 0 stock
	for rows.Next() {
		var productID string
		var newStock int
		if err := rows.Scan(&productID, &newStock); err != nil {
			return fmt.Errorf("error scanning updated stock: %w", err)
		}
		if newStock < 0 {
			return fmt.Errorf("insufficient stock for product ID: %s", productID)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating updated stock: %w", err)
	}

	return nil
}
