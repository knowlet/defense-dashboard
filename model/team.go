package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex"`
	Score int
	Event []Event
}
