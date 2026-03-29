package mysql

import (
	m "gofiber-hax/internal/adapters/db/mysql/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&m.Users{},
	)
}
