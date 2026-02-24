package repository

import (
	"context"
	"database/sql"
	"inventory-backend/internal/core"
	"log"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]core.Category, error)
	Create(ctx context.Context, p core.Category) (core.Category, error)
	Update(ctx context.Context, p core.Category) (core.Category, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (core.Category, error)
}

type postgresCategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &postgresCategoryRepo{db: db}
}

func (p *postgresCategoryRepo) GetAll(ctx context.Context) ([]core.Category, error) {
	query := `
        SELECT 
            id, 
            name, 
            description, 
            created_at
        FROM categories
		LIMIT 30
    `

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Failed in GetAll Categories, Error %v", err)
		return nil, err
	}
	defer rows.Close()

	categories := make([]core.Category, 0)

	for rows.Next() {
		var c core.Category

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.CreatedAt,
		)

		if err != nil {
			log.Printf("Failed in assigning fields, Error %v", err)
			return nil, err
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed in rows.Err(), Error %v", err)
		return nil, err
	}

	return categories, nil
}

func (p *postgresCategoryRepo) GetByID(ctx context.Context, id string) (core.Category, error) {
	var c core.Category
	query := `SELECT 
            id, 
            name, 
            description, 
            created_at
        FROM categories
		WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		log.Printf("Failed in GetByID Categories, Error %v", err)
		return core.Category{}, err
	}

	return c, nil
}
func (p *postgresCategoryRepo) Create(ctx context.Context, c core.Category) (core.Category, error) {
	query := `
        INSERT INTO categories (name, description)
        VALUES ($1, $2)
        RETURNING id, name, description, created_at
    `

	err := p.db.QueryRowContext(ctx, query, c.Name, c.Description).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		log.Printf("Failed in Create Categories, Error %v", err)
		return core.Category{}, err
	}

	return c, nil
}

func (p *postgresCategoryRepo) Update(ctx context.Context, c core.Category) (core.Category, error) {
	query := `
        UPDATE categories
        SET name = $2, description = $3
        WHERE id = $1
        RETURNING id, name, description, created_at
    `

	err := p.db.QueryRowContext(ctx, query, c.ID, c.Name, c.Description).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		log.Printf("Failed in Update Categories, Error %v", err)
		return core.Category{}, err
	}

	return c, nil
}
func (p *postgresCategoryRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Failed in Delete Categories, Error %v", err)
		return err
	}

	return nil
}
