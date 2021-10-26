package score

import (
	"time"

	"defense-dashboard/pkg/helper"
	"defense-dashboard/pkg/quest"

	"gorm.io/gorm"
)

func Scoring(db *gorm.DB, ticker, ticker2 *time.Ticker, quit chan bool) {
	// read data
	q1 := helper.ReadQ1("data/quest1.csv")
	q2 := helper.ReadQ1("data/quest2.csv")
	q3 := helper.ReadQ1("data/quest3.csv")
	q4 := helper.ReadQ1("data/quest4.csv")
	q5 := helper.ReadQ1("data/quest5.csv")
	for {
		select {
		case <-ticker2.C: // check only
			go quest.Subversion(db, q1, true)
			go quest.Exchange(db, q2, true)
			go quest.OA(db, q3, true)
			go quest.Stack(db, q4, true)
			go quest.News(db, q5, true)

		case <-ticker.C:
			// 10/26
			go quest.Subversion(db, q1, false)
			go quest.Exchange(db, q2, false)
			go quest.OA(db, q3, false)
			go quest.Stack(db, q4, false)
			go quest.News(db, q5, false)

		case <-quit:
			ticker.Stop()
		}
	}
}
