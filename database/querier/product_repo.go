package querier

import (
	"context"
	"errors"
	"flukis/invokiss/app/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Querier) Save(ctx context.Context, data model.Product) error {
	query := `
		INSERT INTO products (
			id,
			created_at,
			sku,
			name,
			description,
			image,
			amount
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		) ON CONFLICT(id)
		DO UPDATE SET
			created_at = EXCLUDED.created_at,
			sku = EXCLUDED.sku,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			image = EXCLUDED.image,
			amount = EXCLUDED.amount,
			updated_at = CURRENT_TIMESTAMP;
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

func (q *Querier) Delete(ctx context.Context, data model.Product) error {
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

type ProductWriteModel interface {
	Save(ctx context.Context, data model.Product) error
	Delete(ctx context.Context, data model.Product) error
}

func NewProductWriteModel(
	pool *pgxpool.Pool,
) ProductWriteModel {
	return &Querier{
		pool: pool,
	}
}