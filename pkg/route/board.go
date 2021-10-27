package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// board
func TeamBoardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "board.html", gin.H{})
}

// board/status
func TeamStatusHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "alive.html", gin.H{})
}
