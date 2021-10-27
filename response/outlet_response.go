package response

import (
	"github.com/google/uuid"
	"time"
)

type OutletResponse struct {
	ID          uuid.UUID `json:"id"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
