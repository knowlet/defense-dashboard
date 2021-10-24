package main

import (
	"context"
	"defense-dashboard/model"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

	// Start the server
	r := gin.Default()
	r.GET("/ping", PingHandler)
	r.GET("/service/:status", ServiceHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	select {
	case <-quit: // exit
		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
		return
	}
}
