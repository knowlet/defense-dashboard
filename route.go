package main

import "github.com/gin-gonic/gin"

func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ServiceHandler(c *gin.Context) {
	status := c.Param("status")
	c.JSON(200, gin.H{
		"message": status,
	})
}
