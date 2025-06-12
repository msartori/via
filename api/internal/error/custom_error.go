package custom_error

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewHttpError(status int, msg string) HTTPError {
	return HTTPError{StatusCode: status, Message: msg}
}
