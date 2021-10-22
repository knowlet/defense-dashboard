package model

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Log     string
	Point   int
	TeamID  uint
	QuestID uint
}
