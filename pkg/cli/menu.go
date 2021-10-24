package cli

import (
	"defense-dashboard/model"
	"defense-dashboard/pkg/score"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
)

// TODO: add file lock
var status = false

func svc() string {
	if status {
		return "stop"
	}
	return "start"
}

func Menu(db *gorm.DB, quit chan bool) {
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
					go score.Scoring(db, ticker, stop)
				} else { // stop
					log.Println("Stopping scoring service")
					stop <- true
				}
			}
		case prompt.Options[1]: // view score
			// TODO: add survey list of teams
			queryModel := []model.Team{}
			err := db.Find(&queryModel).Error
			if err != nil {
				fmt.Println(err)
				break
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Score"})

			for _, t := range queryModel {
				type result struct {
					Score int
				}
				var team result
				db.Model(&model.Event{}).
					Select("id, team_id, sum(point) as score").
					Where("team_id = ?", t.ID).
					Find(&team)
				table.Append([]string{t.Name, strconv.Itoa(team.Score)})
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
