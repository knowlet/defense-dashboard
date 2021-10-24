package route

import "github.com/gin-gonic/gin"

func ServiceHandler(c *gin.Context) {
	status := c.Param("status")
	c.JSON(200, gin.H{
		"message": status,
	})
}
