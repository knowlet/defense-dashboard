package quest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

const quest6 = 6

// Chk
func Chk(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			resp, err := request(
				http.MethodGet,
				fmt.Sprintf("http://%s/", t["ip"]),
				t["hostname"].(string),
				nil, nil)
			if err != nil {
				log.Println("[-]", err) // cancel caught
				srvDown(db, quest6, t)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				// read body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("[-]", err)
					srvDown(db, quest6, t)
					return
				}

				// check keywords
				if strings.Contains(string(body), "068349c5c4e75200b9d4cb3a7bb16002") &&
					strings.Contains(string(body), "ca4d0895732e5841bf2f0596bf56e712") &&
					strings.Contains(string(body), "90192f88e905da84e277796cb8a8fc7d") &&
					strings.Contains(string(body), "f3caf0de9e7c164dc18ee9997527feee") &&
					strings.Contains(string(body), "90246a29e6977c1e6ecc7dff6d40e064") {
					plusPoint(db, quest6, t, ischeck)
				} else {
					srvDown(db, quest6, t)
				}
			} else {
				srvDown(db, quest6, t)
			}
		}(team)
	}
}
