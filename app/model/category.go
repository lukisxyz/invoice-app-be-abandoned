package model

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrCategoryNotFound       = errors.New("category: not found")
	ErrCategoryAlreadyDeleted = errors.New("category: already deleted")
)

type Category struct {
	ID        ulid.ULID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCategory(
	Name, Description string,
) Category {
	id := ulid.Make()
	return Category{
		ID:          id,
		CreatedAt:   time.Now(),
		Name:        Name,
		Description: Description,
	}
}
