package response

import (
	"github.com/google/uuid"
	"time"
)

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	OutletID    uuid.UUID `json:"outlet_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stock       int64     `json:"stock"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}
