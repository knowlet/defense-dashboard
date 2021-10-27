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
			h := t["hostname"].(string)
			resp1, err := request(
				http.MethodGet,
				fmt.Sprintf("https://%s/owa/auth/logon.aspx?replaceCurrent=1&url=", t["ip"]),
				h, nil, nil)
			if err != nil {
				log.Println("[-]", err) // cancel caught
				healthcheck(db, quest2, t["id"].(int), ischeck, false)
				return
			}
			defer resp1.Body.Close()
			if resp1.StatusCode != http.StatusOK {
				log.Println("[-]", resp1.StatusCode)
				healthcheck(db, quest2, t["id"].(int), ischeck, false)
				return
			}

			domain := h[strings.Index(h, ".")+1:]
			// querystring
			data := url.Values{}
			data.Set("destination", fmt.Sprintf("https://%s/owa", h))
			data.Set("flags", "4")
			data.Set("forcedownlevel", "0")
			data.Set("username", fmt.Sprintf("JamesHarden@%s", domain))
			data.Set("password", "Def_2021.int")
			data.Set("passwordText", "")
			data.Set("isUtf8", "1")

			log.Println("[+] Querystring:", data.Encode())

			// send login post
			resp, err := request(
				http.MethodPost,
				fmt.Sprintf("https://%s/owa/auth.owa", t["ip"]),
				h,
				strings.NewReader(data.Encode()), nil)
			if err != nil {
				log.Println("[-]", err) // cancel caught
				healthcheck(db, quest2, t["id"].(int), ischeck, false)
				return
			}
			defer resp.Body.Close()
			log.Println("[+]", resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			if resp.StatusCode == http.StatusFound {
				url, err := resp.Location()
				if err != nil {
					log.Println("[-]", err)
					healthcheck(db, quest2, t["id"].(int), ischeck, false)
					return
				}
				log.Println("[+]", url.String())
				healthcheck(db, quest2, t["id"].(int), ischeck, url.Path == "/owa")
				// healthcheck(db, quest2, t["id"].(int), ischeck, url.Path != "")
			} else {
				log.Println("[-]", resp.StatusCode)
				healthcheck(db, quest2, t["id"].(int), ischeck, false)
			}
		}(team)
	}
}
