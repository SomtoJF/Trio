package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
)

// GetCurrentUser godoc
//	@Summary		Get current user
//	@Description	Retrieves the current authenticated user's information
//	@Tags			users
//	@Success		200	{object}	models.User				"Current user data"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/me [get]
func GetCurrentUser(c *gin.Context) {
	user := c.Value("currentUser")
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "We couldn't retrieve your data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUserChats godoc
//	@Summary		Get user chats
//	@Description	Retrieves all chats for the authenticated user
//	@Tags			users
//	@Success		200	{array}		models.Chat				"User chats"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/me/chats [get]
func GetUserChats(c *gin.Context) {
	user := c.Value("currentUser").(models.User)
	var chats []models.Chat

	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chats})
}
