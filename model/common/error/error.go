package error

type Error interface {
	GetMessage() string
	GetStatus() int
}

type NotFoundError struct {
	Message string `json:"message"`
}

func (nfe NotFoundError) GetMessage() string {
	return nfe.Message
}

func (nfe NotFoundError) GetStatus() int {
	return 404
}

type BadRequestError struct {
	Message string `json:"message"`
}

func (bre BadRequestError) GetMessage() string {
	return bre.Message
}

func (bre BadRequestError) GetStatus() int {
	return 400
}
