package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/BMSTU/DIPS/internal/config"
)

// unsetenv снимает переменную окружения на время теста и возвращает её
// после (в testing.T есть Setenv, но нет Unsetenv — поэтому вручную).
func unsetenv(t *testing.T, key string) {
	t.Helper()
	if v, exists := os.LookupEnv(key); exists {
		os.Unsetenv(key)
		t.Cleanup(func() { os.Setenv(key, v) })
	}
}

// TestGetDSN_DatabaseURLHasPriority — Render инжектит DATABASE_URL;
// она должна уходить в GORM без изменений, даже если заданы DB_*.
func TestGetDSN_DatabaseURLHasPriority(t *testing.T) {
	want := "postgres://app:secret@db.example.com:5432/persons?sslmode=require"
	t.Setenv("DATABASE_URL", want)
	t.Setenv("DB_HOST", "localhost") // уступает DATABASE_URL

	if got := config.Load().GetDSN(); got != want {
		t.Errorf("GetDSN() = %q, ожидался DATABASE_URL как есть %q", got, want)
	}
}

// TestGetDSN_FallbackToSeparateVars — локальная разработка: раздельные DB_*.
func TestGetDSN_FallbackToSeparateVars(t *testing.T) {
	t.Setenv("DATABASE_URL", "") // пусто = не задано
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "6543")
	t.Setenv("DB_USER", "ivan")
	t.Setenv("DB_PASSWORD", "pw")
	t.Setenv("DB_NAME", "lab")

	got := config.Load().GetDSN()
	for _, want := range []string{"host=db", "port=6543", "user=ivan", "password=pw", "dbname=lab", "sslmode=disable"} {
		if !strings.Contains(got, want) {
			t.Errorf("GetDSN() = %q, не содержит %q", got, want)
		}
	}
}

// TestServerPort_UsesPlatformPORT — Render задаёт порт через PORT.
func TestServerPort_UsesPlatformPORT(t *testing.T) {
	unsetenv(t, "SERVER_PORT")
	t.Setenv("PORT", "10000")

	if got := config.Load().ServerPort; got != "10000" {
		t.Errorf("ServerPort = %q, ожидался 10000 из переменной PORT", got)
	}
}

// TestServerPort_ServerPortHasPriority — локальная настройка важнее платформенной.
func TestServerPort_ServerPortHasPriority(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("PORT", "10000")

	if got := config.Load().ServerPort; got != "9090" {
		t.Errorf("ServerPort = %q, ожидался 9090 из SERVER_PORT", got)
	}
}
