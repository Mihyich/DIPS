package dto

// ErrorResponse — тело ответа для 404/500.
type ErrorResponse struct {
	Message string `json:"message"`
}

// ValidationErrorResponse — тело ответа для 400.
// Errors: карта «поле → текст ошибки», например {"name": "Name is required"}.
type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
