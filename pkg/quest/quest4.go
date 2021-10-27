package quest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// Stack
func Stack(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	// set loc
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatal(err)
	}
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			// get login page
			resp, err := reqjson(
				http.MethodPost,
				fmt.Sprintf("http://%s/_dash-update-component", t["ip"]),
				t["hostname"].(string),
				strings.NewReader(`{"output":{"id":"graphs_Container","property":"children"},"event":"interval"}`))
			if err != nil {
				log.Println("[-]", err) // cancel caught
				srvDown(db, 4, t)
				return
			}
			defer resp.Body.Close()
			log.Println("[+]", resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			if resp.StatusCode == http.StatusOK {
				// read body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("[-]", err)
					srvDown(db, 4, t)
					return
				}
				json := string(body)
				value := gjson.Get(json, "response.props.children.1.props.figure.layout.title")
				dRexp := regexp.MustCompile(`[\d]{2}:[\d]{2}:[\d]{2}$`)
				// get web time
				dtime := dRexp.FindString(value.String())
				// get current time
				now := time.Now()
				log.Println(now)

				dt := strings.Split(dtime, ":")
				h, _ := strconv.Atoi(dt[0])
				m, _ := strconv.Atoi(dt[1])
				s, _ := strconv.Atoi(dt[2])
				tt := time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
				// deviation within 2 minutes
				sub := now.Sub(tt).Minutes()
				if sub > -2 && sub < 2 {
					plusPoint(db, 4, t, ischeck)
				} else {
					log.Println("[-] Time deviation", sub)
					srvDown(db, 4, t)
				}

				// patch
				// if dtime != "" {
				// 	log.Println("[+] Time", dtime)
				// 	plusPoint(db, 4, t, ischeck)
				// } else {
				// 	srvDown(db, 4, t)
				// }
			} else {
				srvDown(db, 4, t)
			}
		}(team)
	}
}
