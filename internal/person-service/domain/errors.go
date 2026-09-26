package domain

import "errors"

// ErrNotFound — доменная ошибка: сущность не найдена.
// Маппится в HTTP 404 на транспортном слое (handler).
var ErrNotFound = errors.New("person not found")

// ValidationError — ошибка бизнес-валидации с деталями по полям.
// Маппится в HTTP 400 на транспортном слое (handler).
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
