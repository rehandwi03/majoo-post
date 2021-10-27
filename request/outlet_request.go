package request

import "github.com/google/uuid"

type OutletAddRequest struct {
	MerchantID  uuid.UUID `json:"merchant_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Location    string    `json:"location" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required,max=13"`
}

type OutletUpdateRequest struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	MerchantID  uuid.UUID `json:"merchant_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Location    string    `json:"location" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required,max=13"`
}
