package score

import (
	"log"
	"time"

	"defense-dashboard/pkg/helper"
	"defense-dashboard/pkg/quest"

	"gorm.io/gorm"
)

func Scoring(db *gorm.DB, ticker, ticker2 *time.Ticker, quit chan bool) {
	// read data
	// q1 := helper.ReadQ1("data/quest1.csv")
	q2 := helper.ReadQ1("data/quest2.csv")
	// q3 := helper.ReadQ1("data/quest3.csv")
	q4 := helper.ReadQ1("data/quest4.csv")
	// q5 := helper.ReadQ1("data/quest5.csv")
	q6 := helper.ReadQ1("data/quest6.csv")
	q7 := helper.ReadQ1("data/quest7.csv")
	q8 := helper.ReadQ1("data/quest8.csv")
	q9 := helper.ReadQ1("data/quest9.csv")

	for {
		select {
		case <-ticker2.C: // check only
			log.Println("[+] health check start")
			// 10/26
			// go quest.Subversion(db, q1, true)
			go quest.Exchange(db, q2, true)
			// go quest.OA(db, q3, true)
			go quest.Stack(db, q4, true)
			// go quest.News(db, q5, true)

			// 10/28
			go quest.Chk(db, q6, true)
			go quest.Blog(db, q7, true)
			go quest.Git(db, q8, true)
			go quest.Chat(db, q9, true)

		case <-ticker.C:
			log.Println("[+] scoring check start")
			// 10/26
			// go quest.Subversion(db, q1, false)
			go quest.Exchange(db, q2, false)
			// go quest.OA(db, q3, false)
			go quest.Stack(db, q4, false)
			// go quest.News(db, q5, false)

			// 10/28
			go quest.Chk(db, q6, false)
			go quest.Blog(db, q7, false)
			go quest.Git(db, q8, false)
			go quest.Chat(db, q9, false)

		case <-quit:
			ticker.Stop()
			ticker2.Stop()
			return
		}
	}
}
