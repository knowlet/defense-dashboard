package cli

import (
	"defense-dashboard/model"
	"defense-dashboard/pkg/score"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
)

var status = false
var duration = 5 * time.Minute
var duration2 = 45 * time.Second

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
			Options: []string{svc() + " service",
				"view score",
				"lose points",
				"delete logs",
				"exit"},
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
					log.Println("[+] Starting scoring service")
					ticker := time.NewTicker(duration)
					ticker2 := time.NewTicker(duration2)
					go score.Scoring(db, ticker, ticker2, stop)
				} else { // stop
					log.Println("[-] Stopping scoring service")
					stop <- true
				}
			}
		case prompt.Options[1]: // view score
			queryModel := []model.Team{}
			err := db.Find(&queryModel).Error
			if err != nil {
				fmt.Println("[-]", err)
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
			queryModel := []model.Team{}
			err := db.Find(&queryModel).Error
			if err != nil {
				fmt.Println("[-]", err)
				break
			}

			// read quest
			q := []model.Quest{}
			if err := db.Find(&q).Error; err != nil {
				fmt.Println("[-]", err)
				break
			}
			// make quest opts
			qs := []string{}
			for _, q := range queryModel {
				qs = append(qs, q.Name)
			}
			// Choose reason
			reason := ""
			prompt3 := &survey.Input{
				Message: "Why",
				Suggest: func(toComplete string) []string {
					return qs
				},
			}
			survey.AskOne(prompt3, &reason)

			// Choose team
			prompt1 := &survey.MultiSelect{
				Message: "Choose Team:",
				Options: []string{},
			}
			for _, t := range queryModel {
				prompt1.Options = append(prompt1.Options, t.Name)
			}
			teams := []string{}
			survey.AskOne(prompt1, &teams, survey.WithValidator(survey.Required))

			// Choose minus points
			prompt2 := &survey.Select{
				Message: "How many points to minus:",
				Options: []string{"-20", "-30", "-50", "-100", "-20%", "manual input"},
			}
			p := ""
			survey.AskOne(prompt2, &p, survey.WithValidator(survey.Required))

			switch p {
			case prompt2.Options[5]:
				manual := &survey.Input{Message: "How many points to minus:"}
				survey.AskOne(manual, &p)
			}

			for _, team := range teams {
				// read team
				t := model.Team{}
				if err := db.Preload("Events").Where("name = ?", team).Find(&t).Error; err != nil {
					log.Println("[-]", err)
					break
				}

				// check if is persent
				points := 0
				if strings.Contains(p, "%") {
					// get persent
					points, err = strconv.Atoi(p[:strings.Index(p, "%")])
					if err != nil {
						log.Println("[-]", err)
						break
					}
					// sum score
					score := 0
					for _, e := range t.Events {
						score += e.Point
					}
					points = score * points / 100
				} else {
					points, _ = strconv.Atoi(p)
				}
				if err != nil || points == 0 {
					break
				}

				mylog := "[-]"
				if points > 0 {
					mylog = "[+]"
				}
				// save to db
				if err := db.Omit("quest_id").Create(&model.Event{
					Log:    fmt.Sprintf("%s %s %s score %s", mylog, t.Name, reason, p),
					Point:  points,
					TeamID: t.ID,
				}); err != nil {
					log.Println("[-]", err)
				}
				log.Println(mylog, team, points, "reason:", reason)
			}

		case prompt.Options[3]: // delete logs
			queryModel := []model.Event{}
			if err := db.Find(&queryModel).Order("created_at DESC").Error; err != nil {
				log.Println("[-]", err)
				break
			}
			opts := []string{}
			for _, e := range queryModel {
				opts = append(opts, fmt.Sprintf("%d: %s", e.ID, e.Log))
			}
			prompt := &survey.MultiSelect{
				Message: "Choose which log to delete",
				Options: opts,
			}
			logs := []string{}
			survey.AskOne(prompt, &logs)
			// map back to id
			ids := []int{}
			for _, l := range logs {
				id, err := strconv.Atoi(l[:strings.Index(l, ":")])
				if err != nil {
					log.Println("[-]", err)
					continue
				}
				ids = append(ids, id)
			}
			if (ids == nil) || (len(ids) == 0) {
				break
			}
			db.Delete(&model.Event{}, ids)

		default: // exit
			quit <- true
			log.Println("Shutting down server...")
		}
	}

}
