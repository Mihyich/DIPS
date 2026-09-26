package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BMSTU/DIPS/internal/person-service/domain"
	"github.com/BMSTU/DIPS/internal/person-service/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// Мок сервиса. Handler зависит от интерфейса PersonService,
// поэтому в тестах подставляем фейковую реализацию:
// настоящий сервис и БД не нужны.
// ============================================================
type mockPersonService struct {
	persons    []domain.Person
	person     *domain.Person
	listErr    error
	getErr     error
	createErr  error
	updateErr  error
	deleteErr  error
	lastCreate domain.Person
	lastUpdate struct {
		id int64
		p  domain.Person
	}
	lastDeleteID int64
}

func (m *mockPersonService) GetPersons() ([]domain.Person, error) {
	return m.persons, m.listErr
}

func (m *mockPersonService) GetPersonByID(id int64) (*domain.Person, error) {
	return m.person, m.getErr
}

func (m *mockPersonService) CreatePerson(p domain.Person) (*domain.Person, error) {
	m.lastCreate = p
	return m.person, m.createErr
}

func (m *mockPersonService) UpdatePerson(id int64, p domain.Person) (*domain.Person, error) {
	m.lastUpdate.id = id
	m.lastUpdate.p = p
	return m.person, m.updateErr
}

func (m *mockPersonService) DeletePerson(id int64) error {
	m.lastDeleteID = id
	return m.deleteErr
}

func setupRouter(svc PersonService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerRoutes(r, svc) // та же таблица маршрутов, что в проде
	return r
}

// ============================================================
// Тесты
// ============================================================

func TestCreatePerson_Success_Returns201WithLocation(t *testing.T) {
	svc := &mockPersonService{person: &domain.Person{ID: 42, Name: "Ivan"}}
	r := setupRouter(svc)

	body, _ := json.Marshal(dto.PersonRequest{Name: "Ivan", Age: 30})
	req, _ := http.NewRequest("POST", "/api/v1/persons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/api/v1/persons/42", w.Header().Get("Location"))
	assert.Equal(t, "Ivan", svc.lastCreate.Name) // данные дошли до сервиса
}

func TestCreatePerson_Validation_Returns400(t *testing.T) {
	svc := &mockPersonService{
		createErr: &domain.ValidationError{Fields: map[string]string{"name": "Name is required"}},
	}
	r := setupRouter(svc)

	body, _ := json.Marshal(dto.PersonRequest{Name: "Ivan"})
	req, _ := http.NewRequest("POST", "/api/v1/persons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Name is required")
}

func TestGetPersonByID_NotFound_Returns404(t *testing.T) {
	svc := &mockPersonService{getErr: domain.ErrNotFound}
	r := setupRouter(svc)

	req, _ := http.NewRequest("GET", "/api/v1/persons/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPersonByID_ServiceError_Returns500(t *testing.T) {
	svc := &mockPersonService{getErr: assert.AnError} // любая "чужая" ошибка
	r := setupRouter(svc)

	req, _ := http.NewRequest("GET", "/api/v1/persons/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetPersonByID_Success_Returns200(t *testing.T) {
	svc := &mockPersonService{person: &domain.Person{ID: 1, Name: "Peter", Age: 25}}
	r := setupRouter(svc)

	req, _ := http.NewRequest("GET", "/api/v1/persons/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got dto.PersonResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Peter", got.Name)
	assert.Equal(t, 25, got.Age)
}

func TestGetPersons_Empty_ReturnsEmptyJSONArray(t *testing.T) {
	svc := &mockPersonService{persons: nil}
	r := setupRouter(svc)

	req, _ := http.NewRequest("GET", "/api/v1/persons", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String()) // [], а не null
}

func TestUpdatePerson_NotFound_Returns404(t *testing.T) {
	svc := &mockPersonService{updateErr: domain.ErrNotFound}
	r := setupRouter(svc)

	body, _ := json.Marshal(dto.PersonRequest{Name: "Ivan"})
	req, _ := http.NewRequest("PATCH", "/api/v1/persons/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePerson_Success_Returns204(t *testing.T) {
	svc := &mockPersonService{}
	r := setupRouter(svc)

	req, _ := http.NewRequest("DELETE", "/api/v1/persons/7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, int64(7), svc.lastDeleteID) // id дошёл до сервиса
}

func TestDeletePerson_NotFound_Returns404(t *testing.T) {
	svc := &mockPersonService{deleteErr: domain.ErrNotFound}
	r := setupRouter(svc)

	req, _ := http.NewRequest("DELETE", "/api/v1/persons/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
