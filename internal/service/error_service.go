package service

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(msg string) error {
	return &AppError{Code: "BAD_REQUEST", Message: msg}
}

func NotFound(msg string) error {
	return &AppError{Code: "NOT_FOUND", Message: msg}
}

func InvalidInput(msg string) error {
	return &AppError{Code: "INVALID_INPUT", Message: msg}
}

func Conflict(msg string) error {
	return &AppError{Code: "CONFLICT", Message: msg}
}

func Forbidden(msg string) error {
	return &AppError{Code: "FORBIDDEN", Message: msg}
}

func Internal(msg string) error {
	return &AppError{Code: "INTERNAL", Message: msg}
}
