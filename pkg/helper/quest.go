package helper

import (
	"crypto/rand"
	"defense-dashboard/model"
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CsvToMap(file io.Reader) []map[string]interface{} {
	reader := csv.NewReader(file)
	raw, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	csvMap := []map[string]interface{}{}
	header := []string{} // holds first row (header)
	for row, record := range raw {
		// for first row, build the header slice
		if row == 0 {
			for i := 0; i < len(record); i++ {
				header = append(header, strings.TrimSpace(record[i]))
			}
		} else {
			line := map[string]interface{}{}
			for i := 0; i < len(record); i++ {
				line[header[i]] = record[i]
			}
			csvMap = append(csvMap, line)
		}
	}
	return csvMap
}

func SeedQuest(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := CsvToMap(file)
	db.Clauses(clause.OnConflict{DoNothing: true}).
		Model(&model.Quest{}).
		Create(data)
}

func ReadQ1(path string) []map[string]interface{} {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	return CsvToMap(file)
}

func RandomString() string {
	// generate random message
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}
