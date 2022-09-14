package models

import (
	"github.com/qunv/visaurus/core/entities"
	"gorm.io/gorm"
)

func InitializeDatabase(db *gorm.DB) {
	db.Debug().AutoMigrate(
		&entities.Word{},
		&entities.Symnonym{},
	)
}
