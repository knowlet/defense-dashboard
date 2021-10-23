package main

import (
	"defense-dashboard/model"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// points
const (
	plus = 10
)

// TODO: add file lock
var status = false

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
	// TODO: add http server
}
