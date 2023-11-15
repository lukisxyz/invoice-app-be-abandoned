package querier

import (
	"context"
	"flukis/invokiss/app/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *CategoryQuerier) Save(ctx context.Context, data model.Category) error {
	query := `
		INSERT INTO categories (
			id,
			created_at,
			name,
			description
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) ON CONFLICT(id)
		DO UPDATE SET
			created_at = EXCLUDED.created_at,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = CURRENT_TIMESTAMP;
	`
	_, err := q.pool.Exec(
		ctx,
		query,
		data.ID,
		data.CreatedAt,
		data.Name,
		data.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

func (q *CategoryQuerier) Delete(ctx context.Context, data model.Category) error {
	query := `
		UPDATE categories
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

type CategoryWriteModel interface {
	Save(ctx context.Context, data model.Category) error
	Delete(ctx context.Context, data model.Category) error
}

func NewCategoryWriteModel(
	pool *pgxpool.Pool,
) CategoryWriteModel {
	return &CategoryQuerier{
		pool: pool,
	}
}
