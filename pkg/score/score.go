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
	// q2 := helper.ReadQ1("data/quest2.csv")
	// q3 := helper.ReadQ1("data/quest3.csv")
	q4 := helper.ReadQ1("data/quest4.csv")
	// q5 := helper.ReadQ1("data/quest5.csv")
	for {
		select {
		case <-ticker.C:
			// 10/26
			// go quest.Subversion(db, q1)
			// go quest.Exchange(db, q2)
			// go quest.OA(db, q3)
			go quest.Stack(db, q4)
			// go quest.News(db, q5)

		case <-quit:
			ticker.Stop()
		}
	}
}
