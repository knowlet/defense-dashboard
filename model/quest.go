package model

type Quest struct {
	ID     uint   `gorm:"primarykey"`
	Name   string `gorm:"uniqueIndex"`
	Events []Event
}
