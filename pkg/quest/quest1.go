package quest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// Subversion
func Subversion(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			resp, err := request(
				http.MethodGet,
				fmt.Sprintf("http://%s/listing.php", t["ip"]),
				t["hostname"].(string),
				nil)
			if err != nil {
				log.Println(err) // cancel caught
				srvDown(db, 1, t)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				// read body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					srvDown(db, 1, t)
					return
				}

				// check keywords
				if strings.Contains(string(body), "WebSVN") {
					plusPoint(db, 1, t, ischeck)
				} else {
					srvDown(db, 1, t)
				}
			} else {
				srvDown(db, 1, t)
			}
		}(team)
	}
}
