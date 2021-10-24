package seed

import (
	"defense-dashboard/model"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

func SeedQuest(db *gorm.DB, path string) {
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
