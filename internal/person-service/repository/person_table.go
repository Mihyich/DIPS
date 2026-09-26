package repository

import "github.com/BMSTU/DIPS/internal/person-service/domain"

// personTable — сущность уровня БД (таблица persons).
// gorm-теги, имя таблицы и конвертеры живут ТОЛЬКО здесь:
// схема БД не диктует модель бизнеса и не протекает в домен.
type personTable struct {
	ID      int64  `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	Age     int
	Address string
	Work    string
}

// TableName — явно фиксируем имя таблицы (GORM иначе назвал бы её "people").
func (personTable) TableName() string {
	return "persons"
}

// toDomain — конвертер «БД → домен».
func (t personTable) toDomain() domain.Person {
	return domain.Person{
		ID:      t.ID,
		Name:    t.Name,
		Age:     t.Age,
		Address: t.Address,
		Work:    t.Work,
	}
}

// personTableFrom — конвертер «домен → БД».
func personTableFrom(p domain.Person) personTable {
	return personTable{
		ID:      p.ID,
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}
}
