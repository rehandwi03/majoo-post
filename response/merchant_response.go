package response

import (
	"github.com/google/uuid"
	"time"
)

type MerchantResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Name            string    `json:"name"`
	InstitutionName string    `json:"institution_name"`
	PhoneNumber     string    `json:"phone_number"`
	CreatedAt       time.Time `json:"created_at"`
}
