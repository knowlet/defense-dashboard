package quest

import (
	"context"
	"defense-dashboard/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Subversion
func Subversion(db *gorm.DB, data []map[string]interface{}) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/listing.php", t["ip"]), nil)
			req.Host = t["hostname"].(string)
			client := &http.Client{
				// Timeout: 10 * time.Second,
			}
			go func() {
				time.Sleep(time.Second * 10)
				cancel()
			}()
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err) // cancel caught
				return
			}
			if resp.StatusCode == http.StatusOK {
				// read body
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}

				// check keywords
				if strings.Contains(string(body), "WebSVN") {
					// read team info
					team := model.Team{}
					if db.First(&team, t["id"]).Error != nil {
						log.Fatal(err)
					}
					// read quest info
					quest := model.Quest{}
					if db.First(&quest, 1).Error != nil {
						log.Fatal(err)
					}
					// save to db
					db.Create(&model.Event{
						Log:     fmt.Sprintf("%s service alive %s score +%d", quest.Name, team.Name, plus),
						Point:   plus,
						TeamID:  team.ID,
						QuestID: 1,
					})
				}
			}
		}(team)
	}
}
