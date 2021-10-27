package helper

import "github.com/go-playground/validator"

type ErrorResponseValidate struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateRequest(request interface{}) []*ErrorResponseValidate {
	var errors []*ErrorResponseValidate
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponseValidate
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}