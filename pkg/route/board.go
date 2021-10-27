package route

import (
	"defense-dashboard/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// board
func TeamBoardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "board.html", gin.H{})
}

type status struct {
	Tname     string
	Qname     string
	Alive     bool
	CreatedAt time.Time
}

// board/status
func (h Controller) TeamStatusHandler(c *gin.Context) {
	var statuses [][]status
	teams := []model.Team{}
	if err := h.DB.Find(&teams).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	for _, team := range teams {
		var status []status
		subQuery := h.DB.Select("distinct on (quest_id) quest_id, team_id, id, alive, created_at").
			Where("team_id = ?", team.ID).
			Order("quest_id, team_id, id desc").
			Table("statuses")

		if err := h.DB.Table("(?) as s", subQuery).
			Select("teams.name as tname, quests.name as qname, alive, created_at").
			Joins("left join teams on s.team_id = teams.id").
			Joins("left join quests on s.quest_id = quests.id").
			Order("quest_id").
			Find(&status).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		statuses = append(statuses, status)
	}
	c.HTML(http.StatusOK, "alive.html", gin.H{
		"statuses": statuses,
	})
}
