package route

import (
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
	Alive     string
	CreatedAt time.Time
}

// board/status
func (h Controller) TeamStatusHandler(c *gin.Context) {
	var status []status
	subQuery := h.DB.Select("distinct on (quest_id) quest_id, team_id, id, alive, created_at").
		Where("team_id = ?", 1).
		Order("quest_id, team_id, id desc").
		Table("statuses")

	h.DB.Table("(?) as s", subQuery).
		Select("teams.name as tname, quests.name as qname, alive, created_at").
		Joins("left join teams on s.team_id = teams.id").
		Joins("left join quests on s.quest_id = quests.id").
		Order("quest_id").
		Find(&status)

	c.HTML(http.StatusOK, "alive.html", gin.H{
		"status": status,
	})
}
