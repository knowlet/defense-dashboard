package model

import "gorm.io/gorm"

type Quest struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex"`
	Event []Event
}
