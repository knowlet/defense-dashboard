package main

import (
	"time"

	"gorm.io/gorm"
)

func scoring(db *gorm.DB, ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			go quest1(db)
			go quest2(db)

		case <-quit:
			ticker.Stop()
		}
	}

}
