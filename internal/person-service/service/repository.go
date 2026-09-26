package service

import "github.com/BMSTU/DIPS/internal/person-service/domain"

// PersonRepository — порт (интерфейс доступа к данным), который нужен сервису.
// Объявлен здесь, в пакете-потребителе (Go-идиома «интерфейс на стороне
// потребителя»). Репозиторий об этом интерфейсе не знает — он просто
// реализует методы, а main.go сверяет их компилятором.
//
// Отдельный файл внутри пакета service: файлы одного пакета видят друг друга,
// так что это чисто «дробление по концепциям» — порт лежит отдельно от
// бизнес-логики, которая его использует.
type PersonRepository interface {
	FindAll() ([]domain.Person, error)
	FindByID(id int64) (*domain.Person, error)
	Create(person *domain.Person) error
	Update(person *domain.Person) error
	Delete(id int64) error
}
