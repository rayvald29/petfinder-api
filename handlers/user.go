package handlers

import (
	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")

	c.JSON(200, gin.H{
		"message": "User profile retrieved successfully",
		"user_id": userID,
	})

}
