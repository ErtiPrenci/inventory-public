package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-backend/internal/core"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (core.User, error)
}

type postgresUserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) GetByUsername(ctx context.Context, username string) (core.User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	var user core.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.User{}, errors.New("invalid credentials")
		}
		return core.User{}, err
	}
	return user, nil
}
