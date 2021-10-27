package quest

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"defense-dashboard/pkg/helper"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

var (
	pathsChat = []string{
		"/login.aspx",
		"/Registration.aspx",
		"/DownloadFile.aspx",
		"/BankBoard.aspx",
		"/BoardLogin.aspx",
		"/NewsDetail.aspx?ID=2",
	}
)

// Chat
func Chat(db *gorm.DB, data []map[string]interface{}, ischeck bool) {
	for idx, team := range data {
		team["id"] = idx + 1 // team id begins from 1
		go func(t map[string]interface{}) {
			// generate random message
			message := helper.RandomString()
			log.Println("message:", message)
			mac := hmac.New(sha256.New, []byte("HITCON_DEFENSE_2021"))
			mac.Write([]byte(message))
			expectedMAC := mac.Sum(nil)
			log.Println("[+] Expected MAC:", expectedMAC)
			// querystring
			data := url.Values{}
			data.Set("MSG", message)
			log.Println("[+] Querystring:", data.Encode())

			// check all pathsChat
			check := 0
			for _, path := range pathsChat {
				log.Println("[+] Checking", path)

				// create request
				resp, err := request(
					http.MethodPost,
					fmt.Sprintf("http://%s%s", t["ip"], path),
					t["hostname"].(string),
					strings.NewReader(data.Encode()), nil)
				if err != nil {
					log.Println(err) // cancel caught
					srvDown(db, 5, t)
					return
				}
				defer resp.Body.Close()
				log.Println(resp.Request.URL.String())
				log.Println("[+] Response", resp.Status)
				if resp.StatusCode == http.StatusOK {
					// get first line
					scanner := bufio.NewScanner(resp.Body)
					if err != nil {
						srvDown(db, 5, t)
						return
					}
					scanner.Scan()
					line := strings.TrimSpace(scanner.Text())
					log.Println("[+] Response", line)
					// read hmac sha256
					messageMAC, err := hex.DecodeString(line)
					if err != nil {
						srvDown(db, 5, t)
						return
					}
					log.Println("[+] Message MAC:", messageMAC)

					// check hmac
					if hmac.Equal(messageMAC, expectedMAC) {
						check++
						log.Println("[+]", t["ip"], path, "check:", check)
					}
				} else {
					srvDown(db, 5, t)
					return
				}
			}
			if check == len(pathsChat) {
				plusPoint(db, 5, t, ischeck)
			} else {
				srvDown(db, 5, t)
			}
		}(team)
	}
}
