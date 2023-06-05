package handler

const (
	_dbError = 1 + iota
	_apiError
	_validationError
)

type CustomError struct {
	Code    int    `json:"code"`
	Cause   string `json:"cause"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (e *CustomError) Error() string {
	return e.Message
}
