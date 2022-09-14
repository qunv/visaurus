package repositories

import (
	"github.com/qunv/visaurus/core/entities"
	"gorm.io/gorm"
)

type WordRepositories struct {
	db *gorm.DB
}

func NewWordRepositories(db *gorm.DB) *WordRepositories {
	return &WordRepositories{
		db: db,
	}
}

func (w *WordRepositories) Insert(word *entities.Word) error {
	if tx := w.db.Create(word); tx.Error != nil {
		return tx.Error
	}
	return nil
}
