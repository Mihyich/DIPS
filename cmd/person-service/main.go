package main

import (
	"log"

	"github.com/BMSTU/DIPS/internal/config"
	"github.com/BMSTU/DIPS/internal/person-service/repository"
	"github.com/BMSTU/DIPS/internal/person-service/service"
	"github.com/BMSTU/DIPS/internal/person-service/transport/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 1. Инфраструктурная зависимость: подключение к БД.
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 2. Composition root: собираем слои снизу вверх.
	//    Каждый конструктор получает зависимость аргументом (constructor injection).
	repo := repository.NewPersonRepository(db) // repository зависит от db
	svc := service.NewPersonService(repo)      // service зависит от порта PersonRepository
	router := http.NewRouter(svc)              // transport зависит от порта PersonService

	// 3. Запуск. Маршруты зарегистрированы в transport/http/router.go.
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
