package config

import (
	"os"
)

// Конфигурация приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

// Загрузка конфигурации из переменных окружения
func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "program"),
		DBPassword: getEnv("DB_PASSWORD", "test"),
		DBName:     getEnv("DB_NAME", "persons"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

// Формирование строки подключения к базе данных
func (c *Config) GetDSN() string {
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
