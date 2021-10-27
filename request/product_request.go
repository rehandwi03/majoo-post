package request

import (
	"github.com/google/uuid"
)

type ProductAddRequest struct {
	OutletID    uuid.UUID `json:"outlet_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Stock       int64     `json:"stock" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
}

type ProductUpdateRequest struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	OutletID    uuid.UUID `json:"outlet_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Stock       int64     `json:"stock" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
}
