package repository

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-backend/internal/core"
	"log"
)

type ProductRepository interface {
	GetAll(ctx context.Context, filter core.ProductFilter) ([]core.Product, error)
	Count(ctx context.Context, filter core.ProductFilter) (int, error)
	GetByID(ctx context.Context, id string) (core.Product, error)
	Create(ctx context.Context, p core.Product) (core.Product, error)
	Update(ctx context.Context, p core.Product) (core.Product, error)
	Delete(ctx context.Context, id string) error
}

type postgresProductRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepo{db: db}
}

func (r *postgresProductRepo) GetAll(ctx context.Context, filter core.ProductFilter) ([]core.Product, error) {
	query := `
        SELECT 
            id, 
            sku, 
            name, 
            description, 
            category_id, 
            price, 
            cost, 
            current_stock, 
            min_stock, 
            image_url, 
            created_at, 
            updated_at 
        FROM products
        WHERE 1=1
    `

	args := []interface{}{}
	argId := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argId, argId)
		args = append(args, "%"+filter.Search+"%")
		argId++
	}

	if filter.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argId)
		args = append(args, filter.CategoryID)
		argId++
	}

	//Show not deleted products
	query += " AND deleted_at IS NULL"

	// Order by created_at desc by default
	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argId)
		args = append(args, filter.Limit)
		argId++
	}

	if filter.Page > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query += fmt.Sprintf(" OFFSET $%d", argId)
		args = append(args, offset)
		argId++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed in GetAll Products, Error %v", err)
		return nil, err
	}
	defer rows.Close()

	// It starts with make to return [] empty instead of null if there are no data
	products := make([]core.Product, 0)

	for rows.Next() {
		var p core.Product

		err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Description,
			&p.CategoryID,
			&p.Price,
			&p.Cost,
			&p.CurrentStock,
			&p.MinStock,
			&p.ImageURL,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			log.Printf("Failed in assigning fields, Error %v", err)
			return nil, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *postgresProductRepo) Count(ctx context.Context, filter core.ProductFilter) (int, error) {
	query := `SELECT COUNT(*) FROM products WHERE 1=1`
	args := []interface{}{}
	argId := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argId, argId)
		args = append(args, "%"+filter.Search+"%")
		argId++
	}

	if filter.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argId)
		args = append(args, filter.CategoryID)
		argId++
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		log.Printf("Failed in Count Products, Error %v", err)
		return 0, err
	}
	return count, nil
}

func (r *postgresProductRepo) GetByID(ctx context.Context, id string) (core.Product, error) {
	var p core.Product
	query := `SELECT 
            id, 
            sku, 
            name, 
            description, 
            category_id, 
            price, 
            cost, 
            current_stock, 
            min_stock, 
            image_url, 
            created_at, 
            updated_at 
        FROM products
		WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Cost, &p.CurrentStock, &p.MinStock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		log.Printf("Failed in GetByID Products, Error %v", err)
		return core.Product{}, err
	}

	return p, nil
}

func (r *postgresProductRepo) Create(ctx context.Context, p core.Product) (core.Product, error) {
	query := `
        INSERT INTO products (
            sku, 
            name, 
            description, 
            category_id, 
            price, 
            cost, 
            current_stock, 
            min_stock, 
            image_url
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
        RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		p.SKU,
		p.Name,
		p.Description,
		p.CategoryID,
		p.Price,
		p.Cost,
		p.CurrentStock,
		p.MinStock,
		p.ImageURL,
	).Scan(
		&p.ID,
		&p.CreatedAt,
	)

	if err != nil {
		log.Printf("Failed in Create Products, Error %v", err)
		return core.Product{}, err
	}

	return p, nil
}

func (r *postgresProductRepo) Update(ctx context.Context, p core.Product) (core.Product, error) {
	query := `
        UPDATE products
        SET 
            sku = $2, 
            name = $3, 
            description = $4,
            category_id = $5,
            price = $6, 
            cost = $7,
            current_stock = $8, 
            min_stock = $9,
            image_url = $10,
            updated_at = NOW()
        WHERE id = $1::uuid
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(ctx, query,
		p.ID,
		p.SKU,
		p.Name,
		p.Description,
		p.CategoryID,
		p.Price,
		p.Cost,
		p.CurrentStock,
		p.MinStock,
		p.ImageURL,
	).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		log.Printf("Failed in Update Products, Error %v", err)
		return core.Product{}, err
	}

	return p, nil
}

func (r *postgresProductRepo) Delete(ctx context.Context, id string) error {
	//query := `DELETE FROM products WHERE id = $1` //Only by Development
	query := `UPDATE products SET deleted_at = NOW(), updated_at = NOW(), current_stock = 0, min_stock = 0, price = 0, cost = 0 WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Err()
	if err != nil {
		log.Printf("Failed in Delete Products, Error %v", err)
		return err
	}
	return nil
}
