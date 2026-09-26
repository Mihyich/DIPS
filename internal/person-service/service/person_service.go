package service

import "github.com/BMSTU/DIPS/internal/person-service/domain"

// PersonService — слой бизнес-логики. Работает только с доменными типами:
// не знает ни про JSON (dto), ни про GORM (repository).
type PersonService struct {
	repo PersonRepository
}

// NewPersonService — constructor injection: зависимость приходит аргументом.
func NewPersonService(repo PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) GetPersons() ([]domain.Person, error) {
	return s.repo.FindAll()
}

func (s *PersonService) GetPersonByID(id int64) (*domain.Person, error) {
	return s.repo.FindByID(id)
}

func (s *PersonService) CreatePerson(p domain.Person) (*domain.Person, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Create(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdatePerson — частичное обновление (PATCH): меняются только переданные
// поля, неизменённые сохраняются. Полная замена была бы семантикой PUT.
func (s *PersonService) UpdatePerson(id int64, patch domain.PersonPatch) (*domain.Person, error) {
	// Порядок проверок = контракт OpenAPI: сначала 404 (нет записи),
	// потом 400 (нарушено бизнес-правило).
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := patch.Validate(); err != nil {
		return nil, err
	}
	updated := patch.ApplyTo(*existing)
	if err := s.repo.Update(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *PersonService) DeletePerson(id int64) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
