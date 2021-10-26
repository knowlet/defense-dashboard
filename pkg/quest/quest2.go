package quest

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

// Exchange
func Exchange(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			// check logon page
			resp1, err := request(
				http.MethodGet,
				fmt.Sprintf("https://%s/owa/auth/logon.aspx?replaceCurrent=1&url=", t["ip"]),
				t["hostname"].(string), nil)
			if err != nil {
				log.Println(err) // cancel caught
				srvDown(db, 2, t)
				return
			}
			defer resp1.Body.Close()
			if resp1.StatusCode != http.StatusOK {
				srvDown(db, 2, t)
				return
			}

			// querystring
			data := url.Values{}
			data.Set("destination", fmt.Sprintf("https://%s/owa/auth.owa", t["hostname"]))
			data.Set("flags", "4")
			data.Set("forcedownlevel", "0")
			data.Set("username", "JamesHarden@blueteam1.defense.hitcon")
			data.Set("password", "Def_2021.int")
			data.Set("passwordText", "")
			data.Set("isUtf8", "1")

			log.Println("[+] Querystring:", data.Encode())

			// send login post
			resp, err := request(
				http.MethodPost,
				fmt.Sprintf("https://%s/owa/auth.owa", t["ip"]),
				t["hostname"].(string),
				strings.NewReader(data.Encode()))
			if err != nil {
				log.Println(err) // cancel caught
				srvDown(db, 2, t)
				return
			}
			defer resp.Body.Close()
			log.Println(resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			if resp.StatusCode == http.StatusFound {
				url, err := resp.Location()
				if err != nil {
					log.Println(err)
					srvDown(db, 2, t)
					return
				}
				// if url.Path == "/owa" {
				if url.Path != "" {
					plusPoint(db, 2, t, ischeck)
				} else {
					srvDown(db, 2, t)
				}
			} else {
				srvDown(db, 2, t)
			}
		}(team)
	}
}
