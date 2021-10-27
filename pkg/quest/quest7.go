package quest

import (
	"defense-dashboard/pkg/helper"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// Blog
func Blog(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {

			// credentials
			credentials := "healthcheck:U40i nZpF A2I4 sCP9 SoT7 yUgO"
			token := base64.StdEncoding.EncodeToString([]byte(credentials))

			// generate random string
			verfy := helper.RandomString()

			// create post
			resp, err := reqBaseJson(
				http.MethodPost,
				fmt.Sprintf("http://%s/wp-json/wp/v2/posts", t["ip"]),
				t["hostname"].(string),
				token,
				strings.NewReader(fmt.Sprintf(`{
					"title": "health check",
					"content": "%s",
					"status": "publish"
			}`, verfy)))
			if err != nil {
				log.Println(err) // cancel caught
				return
			}
			defer resp.Body.Close()
			log.Println(resp.Request.URL.String())
			log.Println("[+] Response", resp.Status)
			if resp.StatusCode == http.StatusCreated {
				resp2, err := request(
					http.MethodGet,
					fmt.Sprintf("http://%s/", t["ip"]),
					t["hostname"].(string),
					nil, nil)
				if err != nil {
					log.Println(err) // cancel caught
					return
				}
				defer resp2.Body.Close()
				if resp2.StatusCode == http.StatusOK {
					// read body
					body, err := ioutil.ReadAll(resp2.Body)
					if err != nil {
						log.Println(err) // cancel caught
						return
					}
					if strings.Contains(string(body), verfy[:50]) {
						plusPoint(db, 7, t, ischeck)
					} else {
						srvDown(db, 7, t)
					}
				}
			}
		}(team)
	}
}
