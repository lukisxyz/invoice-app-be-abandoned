package querier

import (
	"context"
	"flukis/invokiss/app/model"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type ProductQuerier struct {
	pool *pgxpool.Pool
}

func (q *ProductQuerier) GetOneByID(ctx context.Context, id ulid.ULID) (res model.Product, err error) {
	query := `
		SELECT
			id,
			created_at,
			sku,
			name,
			description,
			image,
			amount,
			updated_at,
			deleted_at
		FROM products
		WHERE id = $1;
	`
	row := q.pool.QueryRow(
		ctx,
		query,
		id,
	)
	var item model.Product
	if err := row.Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Sku,
		&item.Name,
		&item.Description,
		&item.Image,
		&item.Amount,
		&item.UpdatedAt,
		&item.DeletedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return item, model.ErrProductNotFound
		}
		return item, err
	}
	if item.DeletedAt.Valid {
		return item, model.ErrProductAlreadyDeleted
	}
	return item, nil
}

func (q *ProductQuerier) Fetch(ctx context.Context) (res ProductList, err error) {
	var itemCount int

	row := q.pool.QueryRow(
		ctx,
		`
			SELECT
				COUNT(id) as c
			FROM products;

		`,
	)
	if err := row.Scan(&itemCount); err != nil {
		return emptyProducts, err
	}

	if itemCount == 0 {
		return emptyProducts, nil
	}

	items := make([]model.Product, itemCount)
	rows, err := q.pool.Query(
		ctx,
		`
			SELECT
				p.id AS product_id,
				p.sku,
				p.name AS product_name,
				p.description AS product_description,
				p.amount,
				p.image,
				STRING_AGG(c.name, ',' ORDER BY c.name) AS category_names
			FROM
				products p
			LEFT JOIN
				category_product cp ON p.id = cp.product_id
			LEFT JOIN
				categories c ON cp.category_id = c.id
			GROUP BY
				p.id, p.sku, p.name, p.description, p.amount
			ORDER BY p.id;
		`,
	)

	if err != nil {
		return emptyProducts, err
	}

	defer rows.Close()
	var i int
	for i = range items {
		var id ulid.ULID
		var sku string
		var name string
		var description string
		var image *[]byte
		var amount float64
		var cat null.String
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&sku,
			&name,
			&description,
			&amount,
			&image,
			&cat,
		); err != nil {
			return emptyProducts, err
		}

		var categoryNames []string
		if cat.Valid {
			categoryNames = strings.Split(cat.String, ",")
		}

		var categories = make([]model.Category, len(categoryNames))
		for idx := range categoryNames {
			categories[idx].Name = categoryNames[idx]
		}

		items[i] = model.Product{
			ID:          id,
			Sku:         sku,
			Name:        name,
			Description: description,
			Image:       image,
			Amount:      amount,
			Categories:  categories,
		}
	}

	list := ProductList{
		Count: itemCount,
		Data:  items,
	}

	return list, nil
}

type ProductList struct {
	Count int             `json:"count"`
	Data  []model.Product `json:"data"`
}

var emptyProducts = ProductList{
	Count: 0,
	Data:  []model.Product{},
}

type ProductReadModel interface {
	Fetch(ctx context.Context) (res ProductList, err error)
	GetOneByID(ctx context.Context, id ulid.ULID) (res model.Product, err error)
}

func NewProductReadModel(
	pool *pgxpool.Pool,
) ProductReadModel {
	return &ProductQuerier{
		pool: pool,
	}
}
