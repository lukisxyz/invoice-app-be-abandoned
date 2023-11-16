package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrProductSKUDuplicated  = errors.New("product: sku duplicated")
	ErrProductNotFound       = errors.New("product: not found")
	ErrProductAlreadyDeleted = errors.New("product: already deleted")
)

type Product struct {
	ID          ulid.ULID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   null.Time  `json:"updated_at"`
	DeletedAt   null.Time  `json:"deleted_at"`
	Sku         string     `json:"sku"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       *[]byte    `json:"image"`
	Amount      float64    `json:"amount"`
	Categories  []Category `json:"categories"`
}

type CategoryProduct struct {
	CategoryID ulid.ULID `json:"category_id"`
	ProductID  ulid.ULID `json:"product_id"`
}

func NewProduct(
	Sku, Name, Description string,
	Image *[]byte,
	Amount float64,
) Product {
	id := ulid.Make()
	return Product{
		ID:          id,
		CreatedAt:   time.Now(),
		Sku:         Sku,
		Name:        Name,
		Description: Description,
		Image:       Image,
		Amount:      Amount,
	}
}

func (t Product) MarshalJSON() ([]byte, error) {
	s := t.Categories[0].ID.String()
	if s == "00000000000000000000000000" {
		var j struct {
			ID          ulid.ULID `json:"id"`
			CreatedAt   time.Time `json:"created_at"`
			Sku         string    `json:"sku"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Image       *[]byte   `json:"image"`
			Amount      float64   `json:"amount"`
			Categories  []string  `json:"categories"`
		}

		var x = make([]string, len(t.Categories))
		for idx := range t.Categories {
			x[idx] = t.Categories[idx].Name
		}

		j.ID = t.ID
		j.CreatedAt = t.CreatedAt
		j.Sku = t.Sku
		j.Name = t.Name
		j.Description = t.Description
		j.Image = t.Image
		j.Amount = t.Amount
		j.Categories = x

		return json.Marshal(j)
	} else {
		type Cats struct {
			ID   ulid.ULID `json:"id"`
			Name string    `json:"name"`
		}

		var j struct {
			ID          ulid.ULID `json:"id"`
			CreatedAt   time.Time `json:"created_at"`
			Sku         string    `json:"sku"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Image       *[]byte   `json:"image"`
			Amount      float64   `json:"amount"`
			Categories  []Cats    `json:"categories"`
		}

		var x = make([]Cats, len(t.Categories))
		for idx := range t.Categories {
			var bufC = Cats{
				ID:   t.Categories[idx].ID,
				Name: t.Categories[idx].Name,
			}
			x[idx] = bufC
		}

		j.ID = t.ID
		j.CreatedAt = t.CreatedAt
		j.Sku = t.Sku
		j.Name = t.Name
		j.Description = t.Description
		j.Image = t.Image
		j.Amount = t.Amount
		j.Categories = x

		return json.Marshal(j)
	}
}
