package score

import (
	"time"

	"defense-dashboard/pkg/helper"
	"defense-dashboard/pkg/quest"

	"gorm.io/gorm"
)

func Scoring(db *gorm.DB, ticker *time.Ticker, quit chan bool) {
	// read data
	// q1 := helper.ReadQ1("data/quest1.csv")
	q5 := helper.ReadQ1("data/quest5.csv")
	for {
		select {
		case <-ticker.C:
			// 10/26
			// go quest.Subversion(db, q1)
			// go quest.Quest2(db)
			go quest.News(db, q5)

		case <-quit:
			ticker.Stop()
		}
	}
}
