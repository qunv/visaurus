package models

import "gorm.io/gorm"

type Word struct {
	gorm.Model
	Id          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Description string     `json:"description"`
	Symnonyms   []Symnonym `gorm:"foreignKey:WordId;references:Id"`
}

type Symnonym struct {
	gorm.Model
	Id     uint64 `gorm:"primary_key;auto_increment" json:"id"`
	WordId uint64 `json:"word_id"`
}
