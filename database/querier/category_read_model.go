package querier

import (
	"context"
	"flukis/invokiss/app/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type CategoryQuerier struct {
	pool *pgxpool.Pool
}

func (q *CategoryQuerier) Fetch(ctx context.Context) (res CategoryList, err error) {
	var itemCount int

	row := q.pool.QueryRow(
		ctx,
		`
			SELECT
				COUNT(id) as c
			FROM categories;

		`,
	)
	if err := row.Scan(&itemCount); err != nil {
		return emptyCategories, err
	}

	if itemCount == 0 {
		return emptyCategories, nil
	}

	items := make([]model.Category, itemCount)
	rows, err := q.pool.Query(
		ctx,
		`
			SELECT
				id,
				created_at,
				name,
				description,
				updated_at,
				deleted_at
			FROM categories
			ORDER BY id;
		`,
	)

	if err != nil {
		return emptyCategories, err
	}

	defer rows.Close()
	var i int
	for i = range items {
		var id ulid.ULID
		var createdAt time.Time
		var name string
		var description string
		var updateAt null.Time
		var deletedAt null.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&createdAt,
			&name,
			&description,
			&updateAt,
			&deletedAt,
		); err != nil {
			return emptyCategories, err
		}

		items[i] = model.Category{
			ID:          id,
			CreatedAt:   createdAt,
			UpdatedAt:   updateAt,
			DeletedAt:   deletedAt,
			Name:        name,
			Description: description,
		}
	}

	list := CategoryList{
		Count: itemCount,
		Data:  items,
	}

	return list, nil
}

func (q *CategoryQuerier) GetOneByID(ctx context.Context, id ulid.ULID) (res model.Category, err error) {
	query := `
	SELECT
		id,
		created_at,
		name,
		description,
		updated_at,
		deleted_at
	FROM categories
	WHERE id = $1;
`
	row := q.pool.QueryRow(
		ctx,
		query,
		id,
	)
	var item model.Category
	if err := row.Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Name,
		&item.Description,
		&item.UpdatedAt,
		&item.DeletedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return item, model.ErrCategoryNotFound
		}
		return item, err
	}
	if item.DeletedAt.Valid {
		return item, model.ErrCategoryAlreadyDeleted
	}
	return item, nil
}

type CategoryList struct {
	Count int              `json:"count"`
	Data  []model.Category `json:"data"`
}

var emptyCategories = CategoryList{
	Count: 0,
	Data:  []model.Category{},
}

type CategoryReadModel interface {
	Fetch(ctx context.Context) (res CategoryList, err error)
	GetOneByID(ctx context.Context, id ulid.ULID) (res model.Category, err error)
}

func NewCategoryReadModel(
	pool *pgxpool.Pool,
) CategoryReadModel {
	return &CategoryQuerier{
		pool: pool,
	}
}
