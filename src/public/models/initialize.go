package models

import "gorm.io/gorm"

func InitializeDatabase(db *gorm.DB) {
	db.Debug().AutoMigrate(
		&Word{},
		&Symnonym{},
	)
}
