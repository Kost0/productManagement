package repository

import (
	"context"
	"database/sql"
	"products/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	const query = `SELECT * FROM products WHERE id = $1;`

	var product domain.Product

	err := r.db.QueryRowContext(ctx, query, id).Scan()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &product, nil
}
