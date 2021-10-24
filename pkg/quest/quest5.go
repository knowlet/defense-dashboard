package quest

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

var (
	paths = []string{
		"/default.aspx",
		"/Login.aspx",
		"/NewsDetail.aspx?ID=1",
		"/Service.aspx",
		"/DownloadFile.aspx",
	}
)

// News
func News(db *gorm.DB, data []map[string]interface{}) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			check := 0
			for _, path := range paths {
				log.Println("[+] Checking", path)
				// generate random message
				message := make([]byte, 32)
				rand.Read(message)

				// create request
				data := url.Values{}
				data.Set("MSG", string(message))

				resp, err := request(
					http.MethodPost,
					fmt.Sprintf("http://%s%s", t["ip"], path),
					t["hostname"].(string),
					strings.NewReader(data.Encode()))
				if err != nil {
					log.Println(err) // cancel caught
					return
				}
				defer resp.Body.Close()
				log.Println(resp.Request.URL.String())
				log.Println("[+] Response", resp.Status)
				if resp.StatusCode == http.StatusOK {
					// read body
					defer resp.Body.Close()
					// 258d9e21979350a42abb02ad21f60bbf0f398766eea23c126f7ca15acdb33f08
					scanner := bufio.NewScanner(resp.Body)
					if err != nil {
						return
					}
					// get first line
					scanner.Scan()
					messageMAC := []byte(scanner.Text())
					log.Println("[+] Message MAC:", messageMAC)

					// check hmac
					mac := hmac.New(sha256.New, []byte("HITCON_DEFENSE_2021"))
					mac.Write(message)
					expectedMAC := mac.Sum(nil)
					if hmac.Equal(messageMAC, expectedMAC) {
						check++
						log.Println("[+]", t["ip"], path, "check:", check)
					}
				}
			}
			if check == len(paths) {
				plusPoint(db, t)
			}
		}(team)
	}
}
