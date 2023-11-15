package model

import (
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
	ID          ulid.ULID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   null.Time `json:"updated_at"`
	DeletedAt   null.Time `json:"deleted_at"`
	Sku         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"desccription"`
	Image       *[]byte   `json:"image"`
	Amount      float64   `json:"amount"`
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
