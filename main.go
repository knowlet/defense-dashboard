package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"

	"defense-dashboard/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var status = false

func svc() string {
	if status {
		return "stop"
	}
	return "start"
}

func scoring(db *gorm.DB, ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick")
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
				Message: fmt.Sprintf("Are you sure you want to %s the scoring service?", svc()),
				Default: false,
			}
			survey.AskOne(prompt, &check)
			if check {
				status = !status
				if status { // start
					ticker := time.NewTicker(5 * time.Second)
					go scoring(db, ticker, stop)
				} else { // stop
					stop <- true
				}
			}
		case prompt.Options[1]: // view score
		case prompt.Options[2]: // lose points
		case prompt.Options[3]: // exit
			quit <- true
		}
	}
}

func seedTeam(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println()
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
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	raw, err := reader.ReadAll()
	if err != nil {
		panic(err)
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
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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
		return
	}
}
