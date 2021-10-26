package quest

import (
	"defense-dashboard/model"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func plusPoint(db *gorm.DB, qID uint, t map[string]interface{}) {
	// read team info
	team := model.Team{}
	if err := db.First(&team, t["id"]).Error; err != nil {
		log.Fatal(err)
	}
	// read quest info
	quest := model.Quest{}
	if err := db.First(&quest, qID).Error; err != nil {
		log.Fatal(err)
	}
	// save to db
	db.Create(&model.Event{
		Log:     fmt.Sprintf("[+] %s service alive %s score +%d", quest.Name, team.Name, plus),
		Point:   plus,
		TeamID:  team.ID,
		QuestID: qID,
	})
	log.Println(quest.Name, team.Name, plus)
}

func srvDown(db *gorm.DB, qID uint, t map[string]interface{}) {
	// read team info
	team := model.Team{}
	if err := db.First(&team, t["id"]).Error; err != nil {
		log.Fatal(err)
	}
	// read quest info
	quest := model.Quest{}
	if err := db.First(&quest, qID).Error; err != nil {
		log.Fatal(err)
	}
	// save to db
	db.Create(&model.Event{
		Log:     fmt.Sprintf("[-] %s service down %s score +0", quest.Name, team.Name),
		Point:   0,
		TeamID:  team.ID,
		QuestID: qID,
	})
	log.Println("[-]", quest.Name, team.Name, 0)
}
