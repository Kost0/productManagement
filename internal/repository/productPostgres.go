package repository

import (
	"context"
	"database/sql"
	"products/internal/domain"
	"strconv"
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

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.PriceBuy,
		&product.PriceSell,
		&product.SupplierID,
		&product.Weight,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, limit, offset int32) ([]domain.Product, error) {
	const query = `SELECT * FROM products LIMIT $1 OFFSET $2;`

	var products []domain.Product

	rows, err := r.db.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product domain.Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.PriceBuy,
			&product.PriceSell,
			&product.SupplierID,
			&product.Weight,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	const query = `
INSERT INTO products (name, description, priceBuy, priceSell, supplierId, weight)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.PriceBuy,
		product.PriceSell,
		product.SupplierID,
		product.Weight).Scan(&product.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	const query = `
UPDATE products
SET name = $1, description = $2, priceBuy = $3, priceSell = $4, supplierId = $5, weight = $6
WHERE id = $7;
`

	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.PriceBuy,
		product.PriceSell,
		product.SupplierID,
		product.Weight,
		product.ID)

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM products WHERE id = $1;`

	i, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, i)

	if err != nil {
		return err
	}

	return nil
}
