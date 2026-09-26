package domain

// PersonPatch — частичное обновление (PATCH): каждое поле-указатель.
// nil = «поле не передано в запросе» → оставляем существующее значение.
// Указатели — единственный способ в JSON отличить «нет поля» от
// «поле с нулевым значением» ("" / 0): у обычной строки/числа оба
// случая выглядят одинаково — zero value.
type PersonPatch struct {
	Name    *string
	Age     *int
	Address *string
	Work    *string
}

// Validate — правило контракта (OpenAPI: в PersonRequest name обязателен).
// Проверяется на патче ДО слияния: пустое или отсутствующее имя — ошибка 400.
func (pp PersonPatch) Validate() error {
	if pp.Name == nil || *pp.Name == "" {
		return &ValidationError{Fields: map[string]string{"name": "Name is required"}}
	}
	return nil
}

// ApplyTo — сливает переданные поля в существующую сущность.
// Логика живёт в домене: какие поля можно менять и как — бизнес-решение,
// оно одинаково для любого входа (HTTP, CLI, очереди).
func (pp PersonPatch) ApplyTo(p Person) Person {
	if pp.Name != nil {
		p.Name = *pp.Name
	}
	if pp.Age != nil {
		p.Age = *pp.Age
	}
	if pp.Address != nil {
		p.Address = *pp.Address
	}
	if pp.Work != nil {
		p.Work = *pp.Work
	}
	return p
}
