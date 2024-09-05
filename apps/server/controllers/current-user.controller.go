package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
	user := c.Value("currentUser")
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "We couldn't retrieve your data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}
