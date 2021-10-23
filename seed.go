package main

import (
	"bufio"
	"defense-dashboard/model"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

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
