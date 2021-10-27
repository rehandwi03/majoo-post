package helper

type ErrorResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Errors interface{} `json:"errors"`
}

type SuccessResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}


