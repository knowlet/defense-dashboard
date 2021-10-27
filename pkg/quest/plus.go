package quest

import (
	"defense-dashboard/model"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func getInfo(db *gorm.DB, qID, tID uint) (string, string) {
	// read team info
	team := model.Team{}
	if err := db.First(&team, tID).Error; err != nil {
		log.Fatal(err)
	}
	// read quest info
	quest := model.Quest{}
	if err := db.First(&quest, qID).Error; err != nil {
		log.Fatal(err)
	}
	return team.Name, quest.Name
}

func plusPoint(db *gorm.DB, qID, tID uint) {
	tName, qName := getInfo(db, tID, qID)

	// save to db
	db.Create(&model.Event{
		Log:     fmt.Sprintf("[+] %s %s service alive +%d", tName, qName, plus),
		Point:   plus,
		TeamID:  tID,
		QuestID: qID,
	})
	log.Println("[+]", tName, qName, plus)
}

func srvDown(db *gorm.DB, qID, tID uint) {
	tName, qName := getInfo(db, tID, qID)

	// save to db
	db.Create(&model.Event{
		Log:     fmt.Sprintf("[-] %s %s service down +0", tName, qName),
		Point:   0,
		TeamID:  tID,
		QuestID: qID,
	})
	log.Println("[-]", tName, qName, 0)
}

func checkservice(db *gorm.DB, qID, tID uint, alive bool) {
	tName, qName := getInfo(db, tID, qID)

	// save to db
	mylog := ""
	if alive {
		mylog = fmt.Sprintf("[+] %s %s service alive", tName, qName)
	} else {
		mylog = fmt.Sprintf("[-] %s %s service down", tName, qName)
	}

	db.Create(&model.Status{
		Alive:   alive,
		TeamID:  tID,
		QuestID: qID,
	})
	log.Println(mylog)
}

func healthcheck(db *gorm.DB, qID, tID uint, ischeck, isup bool) {
	switch {
	case ischeck:
		checkservice(db, qID, tID, isup)
	case !ischeck && isup:
		plusPoint(db, qID, tID)
	case !ischeck && !isup:
		srvDown(db, qID, tID)
	}
}
