package quest

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
			// generate random message
			key := make([]byte, 32)
			rand.Read(key)
			message := hex.EncodeToString(key)
			log.Println("message:", message)
			mac := hmac.New(sha256.New, []byte("HITCON_DEFENSE_2021"))
			mac.Write([]byte(message))
			expectedMAC := mac.Sum(nil)
			log.Println("[+] Expected MAC:", expectedMAC)
			// querystring
			data := url.Values{}
			data.Set("MSG", message)
			log.Println("[+] Querystring:", data.Encode())

			// check all paths
			check := 0
			for _, path := range paths {
				log.Println("[+] Checking", path)

				// create request
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
					// get first line
					scanner := bufio.NewScanner(resp.Body)
					if err != nil {
						return
					}
					scanner.Scan()
					line := strings.TrimSpace(scanner.Text())
					log.Println("[+] Response", line)
					// read hmac sha256
					messageMAC, err := hex.DecodeString(line)
					if err != nil {
						return
					}
					log.Println("[+] Message MAC:", messageMAC)

					// check hmac
					if hmac.Equal(messageMAC, expectedMAC) {
						check++
						log.Println("[+]", t["ip"], path, "check:", check)
					}
				}
			}
			if check == len(paths) {
				plusPoint(db, 5, t)
			} else {
				srvDown(db, 2, t)
			}
		}(team)
	}
}
