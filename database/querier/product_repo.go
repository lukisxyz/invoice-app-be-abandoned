package querier

import (
	"context"
	"errors"
	"flukis/invokiss/app/model"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

func (q *ProductQuerier) Save(ctx context.Context, data model.Product) error {
	queryInv := `
		INSERT INTO inventories (
			id,
			created_at,
			quantity
		) VALUES (
			$1,
			$2,
			$3
		) ON CONFLICT(id)
		DO UPDATE SET
			created_at = EXCLUDED.created_at,
			quantity = EXCLUDED.quantity,
			updated_at = CURRENT_TIMESTAMP;
	`
	_, err := q.pool.Exec(
		ctx,
		queryInv,
		data.Inventory.ID,
		data.Inventory.CreatedAt,
		data.Inventory.Quantity,
	)

	if err != nil {
		return err
	}

	query := `		
		INSERT INTO products (
			id,
			created_at,
			sku,
			name,
			description,
			image,
			amount,
			inventory_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		) ON CONFLICT(id)
		DO UPDATE SET
			created_at = EXCLUDED.created_at,
			sku = EXCLUDED.sku,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			image = EXCLUDED.image,
			amount = EXCLUDED.amount,
			inventory_id = EXCLUDED.inventory_id,
			updated_at = CURRENT_TIMESTAMP;
	`
	_, err = q.pool.Exec(
		ctx,
		query,
		data.ID,
		data.CreatedAt,
		data.Sku,
		data.Name,
		data.Description,
		data.Image,
		data.Amount,
		data.Inventory.ID,
	)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return model.ErrProductSKUDuplicated
			}
		}

		return err
	}

	return nil
}

func (q *ProductQuerier) Delete(ctx context.Context, data model.Product) error {
	query := `
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`
	_, err := q.pool.Exec(
		ctx,
		query,
		data.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (q *ProductQuerier) AssignCategories(ctx context.Context, productId ulid.ULID, data []ulid.ULID) error {
	var queryIds []string
	var queryValues []any

	for idx := range data {
		queryIds = append(queryIds, fmt.Sprintf("($%d, $%d)", (idx*2+1), (idx*2)+2))
		queryValues = append(queryValues, productId)
		queryValues = append(queryValues, data[idx])
	}

	statement := strings.Join(queryIds, ", ")

	query := fmt.Sprintf(`
		INSERT INTO category_products
			(product_id, category_id)
		VALUES %s
	`, statement)

	_, err := q.pool.Exec(
		ctx,
		query,
		queryValues...,
	)

	if err != nil {
		return err
	}

	return nil
}

// AssignQuantity implements ProductWriteModel.
func (q *ProductQuerier) AssignQuantity(ctx context.Context, productId ulid.ULID, qty int) error {
	query := `
		UPDATE inventories
		SET quantity = $2
		WHERE id = (SELECT inventory_id FROM products WHERE id = $1);
	`
	_, err := q.pool.Exec(
		ctx,
		query,
		productId,
		qty,
	)

	if err != nil {
		return err
	}

	return nil
}

// Edit implements ProductWriteModel.
func (q *ProductQuerier) Edit(ctx context.Context, data model.Product) error {
	query := `		
		UPDATE
			products
		SET
			created_at = $2,
			sku = $3,
			name = $4,
			description = $5,
			image = $6,
			amount = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $1;	
	`
	_, err := q.pool.Exec(
		ctx,
		query,
		data.ID,
		data.CreatedAt,
		data.Sku,
		data.Name,
		data.Description,
		data.Image,
		data.Amount,
	)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return model.ErrProductSKUDuplicated
			}
		}

		return err
	}

	return nil
}

type ProductWriteModel interface {
	Save(ctx context.Context, data model.Product) error
	Edit(ctx context.Context, data model.Product) error
	AssignCategories(ctx context.Context, productId ulid.ULID, data []ulid.ULID) error
	AssignQuantity(ctx context.Context, productId ulid.ULID, qty int) error
	Delete(ctx context.Context, data model.Product) error
}

func NewProductWriteModel(
	pool *pgxpool.Pool,
) ProductWriteModel {
	return &ProductQuerier{
		pool: pool,
	}
}
