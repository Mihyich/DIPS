package repository

import "gorm.io/gorm"

// Migrate создаёт/обновляет таблицы, которыми владеет репозиторий.
// Вызывается один раз из main (composition root).
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&personTable{})
}
