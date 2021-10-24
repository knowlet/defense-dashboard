package score

import (
	"time"

	"defense-dashboard/pkg/quest"

	"gorm.io/gorm"
)

func Scoring(db *gorm.DB, ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			go quest.Quest1(db)
			go quest.Quest2(db)

		case <-quit:
			ticker.Stop()
		}
	}
}
