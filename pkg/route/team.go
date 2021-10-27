package route

import (
	"defense-dashboard/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func (h Controller) TeamHandler(c *gin.Context) {
	id := c.Param("id")
	team := model.Team{}
	if err := h.DB.Preload("Events", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// calc score
	score := 0
	logs := ""
	for _, e := range team.Events {
		score += e.Point
		logs += fmt.Sprintf("[%s] %s\n", e.CreatedAt.Format(time.Kitchen), e.Log)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"team":  team.Name,
		"score": score,
		"logs":  logs,
	})
}

func (h Controller) TeamsHandler(c *gin.Context) {
	queryModel := []model.Team{}
	if err := h.DB.Preload("Events", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Find(&queryModel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	type tt struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
	teams := []tt{}
	for _, t := range queryModel {
		score := 0
		for _, e := range t.Events {
			score += e.Point
		}
		teams = append(teams, tt{
			Name:  t.Name,
			Score: score,
		})
	}
	c.JSON(http.StatusOK, teams)
}

func TeamViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "view.html", gin.H{})
}

func (h Controller) TeamViewLogsHandler(c *gin.Context) {
	nolimit := c.Query("nolimit")
	norefresh := c.Query("norefresh")
	limit := 100
	refresh := true
	if nolimit != "" {
		limit = 0
	}
	if norefresh != "" {
		refresh = false
	}
	queryModel := []struct {
		CreatedAt time.Time
		Log       string
		Point     int
		Name      string
	}{}
	if err := h.DB.Select("*").
		Model(&model.Event{}).
		Joins("left join teams on events.team_id = teams.id").
		Order("created_at DESC").
		Limit(limit).
		Find(&queryModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log events"})
		return
	}

	c.HTML(http.StatusOK, "events.html", gin.H{
		"Events":    queryModel,
		"NoRefresh": refresh,
	})
}

// team/status
func (h Controller) TeamViewStatusHandler(c *gin.Context) {
	nolimit := c.Query("nolimit")
	norefresh := c.Query("norefresh")
	limit := 72
	refresh := true
	if nolimit != "" {
		limit = 0
	}
	if norefresh != "" {
		refresh = false
	}
	queryModel := []struct {
		CreatedAt time.Time
		Alive     bool
		Name      string
		Qname     string
	}{}
	if err := h.DB.Select("statuses.id, alive, team_id, quest_id, teams.name, quests.name as qname, statuses.created_at").
		Model(&model.Status{}).
		Joins("left join teams on statuses.team_id = teams.id").
		Joins("left join quests on statuses.quest_id = quests.id").
		Order("created_at DESC").
		Limit(limit).
		Find(&queryModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log status"})
		return
	}
	c.HTML(http.StatusOK, "status.html", gin.H{
		"Statuses":  queryModel,
		"NoRefresh": refresh,
	})
}

// team/status
func (h Controller) BoardStatusHandler(c *gin.Context) {
	nolimit := c.Query("nolimit")
	norefresh := c.Query("norefresh")
	limit := 72
	refresh := true
	if nolimit != "" {
		limit = 0
	}
	if norefresh != "" {
		refresh = false
	}
	queryModel := []struct {
		CreatedAt time.Time
		Alive     bool
		Name      string
		Qname     string
	}{}
	if err := h.DB.Select("statuses.id, alive, team_id, quest_id, teams.name, quests.name as qname, statuses.created_at").
		Model(&model.Status{}).
		Joins("left join teams on statuses.team_id = teams.id").
		Joins("left join quests on statuses.quest_id = quests.id").
		Order("created_at DESC").
		Limit(limit).
		Find(&queryModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log status"})
		return
	}
	c.HTML(http.StatusOK, "status.html", gin.H{
		"Statuses":  queryModel,
		"NoRefresh": refresh,
	})
}
