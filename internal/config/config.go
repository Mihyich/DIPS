package config

import (
	"os"
)

// Конфигурация приложения
type Config struct {
	// DatabaseURL - подключение в стандартном URL-формате
	// (postgres://user:pass@host:port/dbname). Его автоматически инжектит
	// Render в переменной окружения DATABASE_URL; имеет приоритет
	// перед раздельными DB_*-переменными.
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	// ServerPort - порт HTTP-сервера. Хостинги вроде Render задают
	// стандартную переменную PORT (обычно 10000); SERVER_PORT важнее
	// для локального запуска и docker-compose.
	ServerPort string
}

// Загрузка конфигурации из переменных окружения
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "program"),
		DBPassword:  getEnv("DB_PASSWORD", "test"),
		DBName:      getEnv("DB_NAME", "persons"),
		ServerPort:  getEnv("SERVER_PORT", getEnv("PORT", "8080")),
	}
}

// GetDSN - строка подключения для GORM. Приоритет у DATABASE_URL (Render);
// иначе собирается keyword-DSN из раздельных DB_*-переменных.
// GORM/pgx принимают и URL-формат, и keyword-формат, поэтому URL можно
// отдавать как есть - вместе с ?sslmode=..., если он там есть.
func (c *Config) GetDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable"
}

// Получение значения переменной окружения с возможностью указания значения по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
