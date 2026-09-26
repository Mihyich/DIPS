package dto

import "github.com/BMSTU/DIPS/internal/person-service/domain"

// PersonResponse — представление сущности в ответе API.
// Отдельный тип от domain.Person: поля ответа могут отличаться от домена
// (например, сюда нельзя случайно добавить поле БД или внутренний флаг).
type PersonResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

// NewPersonResponse — явный конвертер «домен → DTO».
// Явные конвертеры вместо автоматических мапперов: видно каждое поле,
// рефакторинг домена ломает компиляцию, а не тихо меняет JSON-контракт.
func NewPersonResponse(p domain.Person) PersonResponse {
	return PersonResponse{
		ID:      p.ID,
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}
}

// NewPersonResponses — конвертер списка; len == 0 даёт [] (JSON-массив),
// а не null: контракт API требует массив даже для пустой базы.
func NewPersonResponses(ps []domain.Person) []PersonResponse {
	out := make([]PersonResponse, 0, len(ps))
	for _, p := range ps {
		out = append(out, NewPersonResponse(p))
	}
	return out
}
