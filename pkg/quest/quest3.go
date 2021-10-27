package quest

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

// OA
func OA(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			// get login page
			resp, err := request(
				http.MethodPost,
				fmt.Sprintf("http://%s/icehrm/app/data/value_Ms7u5RZUJbAv9M1634992053374.png", t["ip"]),
				t["hostname"].(string), nil, nil)
			if err != nil {
				log.Println("[-]", err) // cancel caught
				healthcheck(db, quest3, t["id"].(int), ischeck, false)
				return
			}
			defer resp.Body.Close()
			log.Println("[+]", resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			healthcheck(db, quest3, t["id"].(int), ischeck, resp.StatusCode == http.StatusOK)
		}(team)
	}
}
