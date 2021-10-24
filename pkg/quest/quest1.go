package quest

import (
	"context"
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
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				// read body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}

				// check keywords
				if strings.Contains(string(body), "WebSVN") {
					plusPoint(db, t)
				}
			}
		}(team)
	}
}
