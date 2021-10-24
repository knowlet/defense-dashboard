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
		log.Println()
		if scanner.Text() == "" {
			continue
		}
		queryModel := &model.Team{Name: strings.TrimSpace(scanner.Text())}
		db.FirstOrCreate(&model.Team{}, queryModel)
	}
}
