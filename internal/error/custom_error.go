package custom_error

type NotFoundError struct {
	Message string `json:"message"`
}

func (n *NotFoundError) Error() string {
	return n.Message
}

type ForbiddenError struct {
	Message string `json:"message"`
}

func (f *ForbiddenError) Error() string {
	return f.Message
}

type BadRequest struct {
	Message string `json:"message"`
}

func (b *BadRequest) Error() string {
	return b.Message
}