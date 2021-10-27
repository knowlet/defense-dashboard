package quest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// Git
func Git(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			// check login page
			jar, err := cookiejar.New(nil)
			if err != nil {
				log.Println("[-]", err)
				return
			}
			resp1, err := request(
				http.MethodGet,
				fmt.Sprintf("http://%s/users/sign_in", t["ip"]),
				t["hostname"].(string), nil, jar)
			if err != nil {
				log.Println("[-]", err) // cancel caught
				healthcheck(db, quest8, t["id"].(int), ischeck, false)
				return
			}
			defer resp1.Body.Close()
			if resp1.StatusCode != http.StatusOK {
				healthcheck(db, quest8, t["id"].(int), ischeck, false)
				return
			}

			body, err := io.ReadAll(resp1.Body)
			if err != nil {
				log.Println("[-]", err)
				healthcheck(db, quest8, t["id"].(int), ischeck, false)
				return
			}

			rexp := regexp.MustCompile(`csrf-token" content="(.+)"`)
			// get web time
			sub := rexp.FindStringSubmatch(string(body))
			if len(sub) < 2 {
				log.Println("[-] csrf-token not found")
				healthcheck(db, quest8, t["id"].(int), ischeck, false)
				return
			}
			csrf := sub[1]

			// querystring
			data := url.Values{}
			data.Set("authenticity_token", csrf)
			data.Set("user[login]", "ricky")
			data.Set("user[password]", "Osm3Osm3wwxxd")
			data.Set("user[remember_me]", "0")

			log.Println("[+] Querystring:", data.Encode())

			// send login post
			resp, err := request(
				http.MethodPost,
				fmt.Sprintf("http://%s/users/sign_in", t["ip"]),
				t["hostname"].(string),
				strings.NewReader(data.Encode()), jar)
			if err != nil {
				log.Println("[-]", err) // cancel caught
				healthcheck(db, quest8, t["id"].(int), ischeck, false)
				return
			}
			defer resp.Body.Close()
			log.Println("[+]", resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			if resp.StatusCode == http.StatusFound {
				url, err := resp.Location()
				if err != nil {
					log.Println("[-]", err)
					healthcheck(db, quest8, t["id"].(int), ischeck, false)
					return
				}
				healthcheck(db, quest8, t["id"].(int), ischeck, url.Path == "/")
			} else {
				healthcheck(db, quest8, t["id"].(int), ischeck, false)
			}
		}(team)
	}
}
