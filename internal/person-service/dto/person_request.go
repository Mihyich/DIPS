package dto

import "github.com/BMSTU/DIPS/internal/person-service/domain"

// PersonRequest — тело HTTP-запроса для POST/PATCH /api/v1/persons.
// Только JSON-теги: описывает контракт API, а не схему БД.
// Нет поля ID: клиент не может назначить ID записи сам.
type PersonRequest struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

// ToDomain — переводит DTO в доменную сущность.
// Граница «транспорт → домен»: JSON-формат исчезает, остаётся бизнес-тип.
func (r PersonRequest) ToDomain() domain.Person {
	return domain.Person{
		Name:    r.Name,
		Age:     r.Age,
		Address: r.Address,
		Work:    r.Work,
	}
}

// PersonPatchRequest — тело PATCH: все поля — указатели, потому что
// частичное обновление обязано отличать «нет поля» от «пустое значение».
// nil = в JSON поля не было → поле не трогаем. Имеется в виду только то,
// что клиент явно прислал.
type PersonPatchRequest struct {
	Name    *string `json:"name"`
	Age     *int    `json:"age"`
	Address *string `json:"address"`
	Work    *string `json:"work"`
}

// ToDomain — переводит DTO в доменный патч (граница «транспорт → домен»).
func (r PersonPatchRequest) ToDomain() domain.PersonPatch {
	return domain.PersonPatch{
		Name:    r.Name,
		Age:     r.Age,
		Address: r.Address,
		Work:    r.Work,
	}
}
