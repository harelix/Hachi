package api

type APIError struct {
	msg string
}

func (error *APIError) Error() string {
	return error.msg
}
func APIErrorResponse(message string) error {
	return &APIError{message}
}
