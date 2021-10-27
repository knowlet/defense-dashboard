package model

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Alive   bool
	TeamID  uint
	QuestID uint
}
