package quest

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

// OA
func OA(db *gorm.DB, data []map[string]interface{}) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			// get login page
			resp, err := request(
				http.MethodPost,
				fmt.Sprintf("https://%s/icehrm/app/login.php", t["ip"]),
				t["hostname"].(string), nil)
			if err != nil {
				log.Println(err) // cancel caught
				return
			}
			defer resp.Body.Close()
			log.Println(resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			if resp.StatusCode == http.StatusOK {
				plusPoint(db, 3, t)
			}
		}(team)
	}
}
