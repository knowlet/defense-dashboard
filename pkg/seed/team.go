package seed

import (
	"bufio"
	"defense-dashboard/model"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

func SeedTeam(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		teamName := strings.TrimSpace(scanner.Text())
		if teamName == "" {
			continue
		}
		queryModel := &model.Team{Name: teamName}
		db.FirstOrCreate(&model.Team{}, queryModel)
	}
}
