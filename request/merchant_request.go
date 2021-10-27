package request

import "github.com/google/uuid"

type MerchantAddRequest struct {
	Name            string `json:"name" validate:"required,max=255"`
	InstitutionName string `json:"institution_name" validate:"required,max=255"`
	PhoneNumber     string `json:"phone_number" validate:"required,max=13"`
}

type MerchantUpdateRequest struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	Name            string    `json:"name" validate:"required,max=255"`
	InstitutionName string    `json:"institution_name" validate:"required,max=255"`
	PhoneNumber     string    `json:"phone_number" validate:"required,max=13"`
}
