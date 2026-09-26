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

func (s *PersonService) UpdatePerson(id int64, p domain.Person) (*domain.Person, error) {
	// Порядок проверок = контракт OpenAPI: сначала 404 (нет записи),
	// потом 400 (нарушено бизнес-правило).
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	existing.Name = p.Name
	existing.Age = p.Age
	existing.Address = p.Address
	existing.Work = p.Work
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *PersonService) DeletePerson(id int64) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
