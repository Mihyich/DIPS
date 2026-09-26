package repository

import (
	"errors"

	"github.com/BMSTU/DIPS/internal/person-service/domain"
	"gorm.io/gorm"
)

// GormPersonRepository — реализация service.PersonRepository поверх GORM
// (adaptor в терминах hexagonal architecture). Единственный пакет,
// который знает и про GORM, и про домен.
type GormPersonRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

func (r *GormPersonRepository) FindAll() ([]domain.Person, error) {
	var tables []personTable

	if err := r.db.Find(&tables).Error; err != nil {
		return nil, err
	}

	persons := make([]domain.Person, 0, len(tables))

	for _, t := range tables {
		persons = append(persons, t.toDomain())
	}

	return persons, nil
}

func (r *GormPersonRepository) FindByID(id int64) (*domain.Person, error) {
	var t personTable

	err := r.db.First(&t, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Переводим ошибку GORM в доменную: верхние слои не знают про gorm.
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	person := t.toDomain()

	return &person, nil
}

func (r *GormPersonRepository) Create(person *domain.Person) error {
	t := personTableFrom(*person)

	if err := r.db.Create(&t).Error; err != nil {
		return err
	}

	person.ID = t.ID // GORM записал сгенерированный ID обратно

	return nil
}

func (r *GormPersonRepository) Update(person *domain.Person) error {
	return r.db.Save(personTableFrom(*person)).Error
}

func (r *GormPersonRepository) Delete(id int64) error {
	return r.db.Delete(&personTable{}, id).Error
}
