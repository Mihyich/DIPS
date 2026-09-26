package domain

// Person — доменная сущность: центральное понятие бизнеса.
// Чистый Go-тип: НЕТ JSON-тегов, НЕТ gorm-тегов, НЕТ HTTP.
// Домен ничего не знает о формате API, схеме БД и веб-фреймворке.
type Person struct {
	ID      int64
	Name    string
	Age     int
	Address string
	Work    string
}

// Validate — бизнес-правило домена: имя обязательно.
// Правило живёт здесь, а не в handler, — оно одинаково для любого входа
// (HTTP, CLI, очереди и т.д.).
func (p Person) Validate() error {
	if p.Name == "" {
		return &ValidationError{Fields: map[string]string{"name": "Name is required"}}
	}
	return nil
}
