package http

import "github.com/gin-gonic/gin"

// NewRouter собирает движок Gin со всеми маршрутами сервиса.
// Транспортный уровень сам владеет своим HTTP-контрактом:
// main.go только вызывает NewRouter и передаёт ему собранную зависимость.
func NewRouter(svc PersonService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil) // Прод-логи дополнительно чистятся переменной GIN_MODE=release.
	registerRoutes(r, svc)
	return r
}

// registerRoutes — таблица маршрутов. Используется и в проде (через NewRouter),
// и в тестах напрямую: тесты берут gin.New() без логгера, чтобы вывод был чистым.
func registerRoutes(r *gin.Engine, svc PersonService) {
	h := NewPersonHandler(svc)
	api := r.Group("/api/v1")
	{
		api.GET("/persons", h.GetPersons)
		api.GET("/persons/:id", h.GetPersonByID)
		api.POST("/persons", h.CreatePerson)
		api.PATCH("/persons/:id", h.UpdatePerson)
		api.DELETE("/persons/:id", h.DeletePerson)
	}
}
