package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
)

func GetCurrentUser(c *gin.Context) {
	user := c.Value("currentUser")
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "We couldn't retrieve your data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUserChats(c *gin.Context) {
	user := c.Value("currentUser").(models.User)
	var chats []models.Chat

	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chats})
}
