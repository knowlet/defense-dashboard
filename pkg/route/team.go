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

func TeamViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "view.html", gin.H{})
}
