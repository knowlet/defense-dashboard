package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"

	"defense-dashboard/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var status = false

func svc() string {
	if status {
		return "stop"
	}
	return "start"
}

// points
const (
	plus = 10
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

// quest2
func quest2(db *gorm.DB) {
	log.Println("quest2")
}

func scoring(db *gorm.DB, ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			go quest1(db)
			go quest2(db)

		case <-quit:
			ticker.Stop()
		}
	}

}

func menu(db *gorm.DB, quit chan bool) {
	stop := make(chan bool)
	for {
		prompt := &survey.Select{
			Message: "Welcome to Dashboard System:",
			Options: []string{svc() + " service", "view score", "lose points", "exit"},
		}
		var opts string
		survey.AskOne(prompt, &opts, survey.WithValidator(survey.Required))
		switch opts {
		case prompt.Options[0]: // start/stop service
			check := false
			prompt := &survey.Confirm{
				Message: "Are you sure you want to " + svc() + " the scoring service?",
				Default: true,
			}
			survey.AskOne(prompt, &check)
			if check {
				status = !status
				if status { // start
					log.Println("Starting scoring service")
					ticker := time.NewTicker(5 * time.Second)
					go scoring(db, ticker, stop)
				} else { // stop
					log.Println("Stopping scoring service")
					stop <- true
				}
			}
		case prompt.Options[1]: // view score
			queryModel := []model.Team{}
			err := db.Preload(clause.Associations).Find(&queryModel).Error
			if err != nil {
				fmt.Println(err)
				break
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Score"})

			for _, v := range queryModel {
				// TODO: use sql sum
				sum := 0
				for _, s := range v.Events {
					sum += s.Point
				}
				table.Append([]string{v.Name, strconv.Itoa(sum)})
			}
			table.Render() // Send output
		case prompt.Options[2]: // lose points
		default: // exit
			quit <- true
		}
	}
}

func seedTeam(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.Println()
		if scanner.Text() == "" {
			continue
		}
		queryModel := &model.Team{Name: strings.TrimSpace(scanner.Text())}
		db.FirstOrCreate(&model.Team{}, queryModel)
	}
}

func seedQuest(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	raw, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	header := []string{} // holds first row (header)
	for row, record := range raw {
		// for first row, build the header slice
		if row == 0 {
			for i := 0; i < len(record); i++ {
				header = append(header, strings.TrimSpace(record[i]))
			}
		} else {
			line := map[string]string{}
			for i := 0; i < len(record); i++ {
				line[header[i]] = record[i]
			}
			db.FirstOrCreate(&model.Quest{}, line)
		}
	}
}

func main() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("verbose.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatal(err)
	}
	// Set log to file
	log.SetOutput(file)

	// Open the data.db file. It will be created if it doesn't exist.
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		model.Team{},
		model.Quest{},
		model.Event{},
	)

	// Seed team data from file
	seedTeam(db, "data/teams.txt")
	// Seed quest data from file
	seedQuest(db, "data/quests.csv")

	// Start the menu
	quit := make(chan bool)
	go menu(db, quit)

	select {
	case <-quit: // exit
		log.Println("Bye")
		return
	}
}
