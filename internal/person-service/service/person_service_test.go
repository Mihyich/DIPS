package service

import (
	"testing"

	"github.com/BMSTU/DIPS/internal/person-service/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// Мок репозитория. Service зависит от интерфейса PersonRepository,
// поэтому тесты работают без БД и без GORM.
// ============================================================
type mockPersonRepository struct {
	persons    map[int64]*domain.Person
	nextID     int64
	createErr  error
	findErr    error
	deleteErr  error
	findAllErr error
	created    *domain.Person
	deletedID  int64
}

func newMockPersonRepository() *mockPersonRepository {
	return &mockPersonRepository{
		persons: map[int64]*domain.Person{},
		nextID:  1,
	}
}

func (m *mockPersonRepository) FindAll() ([]domain.Person, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	out := make([]domain.Person, 0, len(m.persons))
	for _, p := range m.persons {
		out = append(out, *p)
	}
	return out, nil
}

func (m *mockPersonRepository) FindByID(id int64) (*domain.Person, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if p, ok := m.persons[id]; ok {
		return p, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockPersonRepository) Create(person *domain.Person) error {
	if m.createErr != nil {
		return m.createErr
	}
	person.ID = m.nextID
	m.nextID++
	m.persons[person.ID] = person
	m.created = person
	return nil
}

func (m *mockPersonRepository) Update(person *domain.Person) error {
	m.persons[person.ID] = person
	return nil
}

func (m *mockPersonRepository) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deletedID = id
	delete(m.persons, id)
	return nil
}

// ============================================================
// Тесты бизнес-логики сервиса
// ============================================================

// strPtr / intPtr — удобители для полей-указателей в domain.PersonPatch:
// nil означает «поле не передано», непустой указатель — «обновить».
func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

func TestCreatePerson_ValidRequest_PersistsAndReturnsID(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)

	person, err := svc.CreatePerson(domain.Person{Name: "Ivan", Age: 30, Work: "BMSTU"})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), person.ID) // ID назначен "БД"
	assert.Equal(t, "Ivan", repo.created.Name)
}

func TestCreatePerson_EmptyName_ReturnsValidationError(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)

	person, err := svc.CreatePerson(domain.Person{Name: "", Age: 30})

	assert.Nil(t, person)
	assert.Error(t, err)

	var verr *domain.ValidationError
	assert.ErrorAs(t, err, &verr) // ошибка именно валидации
	assert.Equal(t, "Name is required", verr.Fields["name"])
	assert.Nil(t, repo.created) // в "БД" ничего не ушло
}

func TestGetPersonByID_Missing_ReturnsErrNotFound(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)

	person, err := svc.GetPersonByID(999)

	assert.Nil(t, person)
	assert.ErrorIs(t, err, domain.ErrNotFound) // доменная ошибка, не gorm
}

func TestUpdatePerson_Missing_ReturnsErrNotFound(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)

	person, err := svc.UpdatePerson(999, domain.PersonPatch{Name: strPtr("Ivan")})

	assert.Nil(t, person)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUpdatePerson_Existing_ChangesFields(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)
	_, _ = svc.CreatePerson(domain.Person{Name: "Ivan", Age: 30})

	updated, err := svc.UpdatePerson(1, domain.PersonPatch{Name: strPtr("Ivan Jr"), Age: intPtr(31)})

	assert.NoError(t, err)
	assert.Equal(t, "Ivan Jr", updated.Name)
	assert.Equal(t, 31, updated.Age)
	assert.Equal(t, int64(1), updated.ID) // ID не изменился
}

// TestUpdatePerson_PartialPatch_PreservesUnspecifiedFields — регрессия на баг
// интеграционных тестов: PATCH приходит с полями только name и address,
// а work/age обязаны сохраниться (раньше затирались нулевыми значениями).
func TestUpdatePerson_PartialPatch_PreservesUnspecifiedFields(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)
	_, _ = svc.CreatePerson(domain.Person{Name: "Ivan", Age: 31, Address: "Moscow", Work: "BMSTU"})

	updated, err := svc.UpdatePerson(1, domain.PersonPatch{
		Name:    strPtr("Petr"),
		Address: strPtr("Tula"),
	})

	assert.NoError(t, err)
	assert.Equal(t, "Petr", updated.Name)    // передано → изменено
	assert.Equal(t, "Tula", updated.Address) // передано → изменено
	assert.Equal(t, 31, updated.Age)         // НЕ передано → сохранилось (не 0)
	assert.Equal(t, "BMSTU", updated.Work)   // НЕ передано → сохранилось (не "")
	assert.Equal(t, int64(1), updated.ID)
}

func TestUpdatePerson_InvalidName_ReturnsValidationError(t *testing.T) {
	cases := map[string]domain.PersonPatch{
		"no name field": {},                 // спека: name обязателен
		"empty name":    {Name: strPtr("")}, // бизнес-правило домена
	}
	for name, patch := range cases {
		t.Run(name, func(t *testing.T) {
			repo := newMockPersonRepository()
			svc := NewPersonService(repo)
			_, _ = svc.CreatePerson(domain.Person{Name: "Ivan"})

			person, err := svc.UpdatePerson(1, patch)

			assert.Nil(t, person)
			var verr *domain.ValidationError
			assert.ErrorAs(t, err, &verr) // ошибка именно валидации
			assert.Equal(t, "Name is required", verr.Fields["name"])
			assert.Equal(t, "Ivan", repo.persons[1].Name) // в "БД" ничего не изменилось
		})
	}
}

func TestDeletePerson_Missing_ReturnsErrNotFound(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)

	err := svc.DeletePerson(999)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Equal(t, int64(0), repo.deletedID) // Delete не вызывался
}

func TestDeletePerson_Existing_RemovesFromRepo(t *testing.T) {
	repo := newMockPersonRepository()
	svc := NewPersonService(repo)
	_, _ = svc.CreatePerson(domain.Person{Name: "Ivan"})

	err := svc.DeletePerson(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), repo.deletedID)
}
