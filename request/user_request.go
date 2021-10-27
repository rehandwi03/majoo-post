package request

import (
	"github.com/google/uuid"
)

type UserAddRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=3,max=50"`
	LastName    string `json:"last_name" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,min=3,max=13"`
}

type UserUpdateRequest struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	FirstName   string    `json:"first_name" validate:"required,min=3,max=50"`
	LastName    string    `json:"last_name" validate:"required,min=3,max=50"`
	Email       string    `json:"email" validate:"required,email,min=3,max=50"`
	Password    string    `json:"password" validate:"required,min=8,max=50"`
	PhoneNumber string    `json:"phone_number" validate:"required,min=3,max=13"`
}
