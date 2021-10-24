package main

import (
	"defense-dashboard/model"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func svc() string {
	if status {
		return "stop"
	}
	return "start"
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
			// TODO: add survey list of teams
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
			// TODO: add lose points events
		default: // exit
			quit <- true
			log.Println("Shutting down server...")
		}
	}

}
