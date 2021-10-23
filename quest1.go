package main

import (
	"context"
	"defense-dashboard/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// quest1
func quest1(db *gorm.DB) {
	type t struct {
		id       uint
		ip       string
		hostname string
		pass     bool
	}

	var teams = []t{
		{1, "127.0.0.1", "example.com", false},
	}

	for _, team := range teams {
		go func(team t) {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:2021", team.ip), nil)
			req.Host = team.hostname
			client := &http.Client{}
			// Timeout: 5 * time.Second,
			go func() {
				time.Sleep(time.Second * 5)
				cancel()
			}()
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err) // cancel caught
				return
			}
			log.Println(resp.StatusCode)
			// save to db
			db.Create(&model.Event{
				Log:     fmt.Sprintf("#%d: Service alive Team%d score +%d", 1, team.id, plus),
				Point:   plus,
				TeamID:  team.id,
				QuestID: 1,
			})
		}(team)
	}
}
