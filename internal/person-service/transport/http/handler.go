package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BMSTU/DIPS/internal/person-service/domain"
	"github.com/BMSTU/DIPS/internal/person-service/dto"
	"github.com/gin-gonic/gin"
)

// PersonService — контракт, который нужен HTTP-слою.
// Объявлен здесь, в пакете-потребителе. Handler работает с доменом и dto
// и не знает, какой сервис за ним стоит (настоящий, мок, фейк).
type PersonService interface {
	GetPersons() ([]domain.Person, error)
	GetPersonByID(id int64) (*domain.Person, error)
	CreatePerson(p domain.Person) (*domain.Person, error)
	UpdatePerson(id int64, p domain.Person) (*domain.Person, error)
	DeletePerson(id int64) error
}

// PersonHandler — транспортный слой: парсит HTTP, вызывает сервис,
// маппит доменные ошибки в HTTP-статусы. Бизнес-логики здесь нет.
type PersonHandler struct {
	svc PersonService
}

// NewPersonHandler — constructor injection.
func NewPersonHandler(svc PersonService) *PersonHandler {
	return &PersonHandler{svc: svc}
}

// writeError — единая точка маппинга «доменная ошибка → HTTP».
func writeError(c *gin.Context, err error) {
	var verr *domain.ValidationError
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Person not found"})
	case errors.As(err, &verr):
		c.JSON(http.StatusBadRequest, dto.ValidationErrorResponse{
			Message: "Validation failed",
			Errors:  verr.Fields,
		})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Internal server error"})
	}
}

// GetPersons godoc
// @Summary Get all Persons
// @Produce json
// @Success 200 {array} dto.PersonResponse
// @Router /api/v1/persons [get]
func (h *PersonHandler) GetPersons(c *gin.Context) {
	persons, err := h.svc.GetPersons()
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPersonResponses(persons))
}

// GetPersonByID godoc
// @Summary Get Person by ID
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} dto.PersonResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/persons/{id} [get]
func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid person ID"})
		return
	}

	person, err := h.svc.GetPersonByID(id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPersonResponse(*person))
}

// CreatePerson godoc
// @Summary Create new Person
// @Accept json
// @Produce json
// @Param person body dto.PersonRequest true "Person data"
// @Success 201 {string} string "Location header"
// @Failure 400 {object} dto.ValidationErrorResponse
// @Router /api/v1/persons [post]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationErrorResponse{
			Message: "Validation failed",
			Errors:  map[string]string{"body": "Invalid JSON"},
		})
		return
	}

	person, err := h.svc.CreatePerson(req.ToDomain())
	if err != nil {
		writeError(c, err)
		return
	}

	// Требование ТЗ: 201 Created + заголовок Location.
	c.Header("Location", "/api/v1/persons/"+strconv.FormatInt(person.ID, 10))
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdatePerson godoc
// @Summary Update Person by ID
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body dto.PersonRequest true "Person data"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {object} dto.ValidationErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/persons/{id} [patch]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid person ID"})
		return
	}

	var req dto.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ValidationErrorResponse{
			Message: "Validation failed",
			Errors:  map[string]string{"body": "Invalid JSON"},
		})
		return
	}

	person, err := h.svc.UpdatePerson(id, req.ToDomain())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPersonResponse(*person))
}

// DeletePerson godoc
// @Summary Remove Person by ID
// @Produce json
// @Param id path int true "Person ID"
// @Success 204 {string} string ""
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/persons/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid person ID"})
		return
	}

	if err := h.svc.DeletePerson(id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
