package criteria

import "github.com/rehandwi03/test-case-backend-majoo/util"

type UserCriteria struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Pagination  util.Pagination
}
