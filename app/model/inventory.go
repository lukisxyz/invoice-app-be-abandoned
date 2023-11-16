package model

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrInventoryNotFound       = errors.New("inventory: not found")
	ErrInventoryAlreadyDeleted = errors.New("inventory: already deleted")
)

type Inventory struct {
	ID        ulid.ULID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`

	Quantity int `json:"quantity"`
}

func NewInventory(
	Quantity int, Description string,
) Inventory {
	id := ulid.Make()
	return Inventory{
		ID:        id,
		CreatedAt: time.Now(),
		Quantity:  Quantity,
	}
}
