package model

type Team struct {
	ID     uint   `gorm:"primarykey"`
	Name   string `gorm:"uniqueIndex"`
	Score  int
	Events []Event
}
