package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
	"gorm.io/gorm"
)

// CreateChat creates a new chat for the authenticated user
func CreateChat(c *gin.Context) {
	var body struct {
		ChatName string `json:"chatName" binding:"required,max=20"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	chat := models.Chat{
		ChatName: body.ChatName,
		UserID:   userModel.ID,
	}

	result := initializers.DB.Create(&chat)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": chat})
}

// AddAgentToChat adds an agent to a chat (max 2 agents per chat)
func AddAgentToChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body struct {
		Name   string   `json:"name" binding:"required,max=20"`
		Lingo  string   `json:"lingo" binding:"required,max=20"`
		Traits []string `json:"traits" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.Traits) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent cannot have more than 4 traits"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var chat models.Chat
	if err := initializers.DB.Preload("Agents").First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	if len(chat.Agents) >= 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat already has the maximum number of agents"})
		return
	}

	// Check if an agent with the same name already exists in this chat
	for _, existingAgent := range chat.Agents {
		if existingAgent.Name == body.Name {
			c.JSON(http.StatusConflict, gin.H{"error": "An agent with this name already exists in the chat"})
			return
		}
	}

	agent := models.Agent{
		Name:   body.Name,
		Lingo:  body.Lingo,
		Traits: body.Traits,
		ChatID: chat.ID,
	}

	if err := initializers.DB.Create(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": agent})
}

// DeleteChat deletes a chat
func DeleteChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var chat models.Chat
	if err := initializers.DB.First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	if err := initializers.DB.Delete(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Chat deleted successfully"})
}

// UpdateChat updates a chat's name
func UpdateChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body struct {
		ChatName string `json:"chatName" binding:"required,max=20"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var chat models.Chat
	if err := initializers.DB.First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	chat.ChatName = body.ChatName

	if err := initializers.DB.Save(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chat})
}

// GetChatInfo retrieves chat information including its agents and messages with sender details
func GetChatInfo(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var chat models.Chat
	if err := initializers.DB.Preload("Agents").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Prepare the response
	type MessageWithSender struct {
		models.Message
		Sender interface{} `json:"sender"`
	}

	var messagesWithSenders []MessageWithSender

	for _, message := range chat.Messages {
		var sender interface{}

		if message.SenderType == string(types.SenderTypeUser) {
			var user models.User
			initializers.DB.First(&user, message.SenderID)
			sender = struct {
				ID       uuid.UUID `json:"id"`
				Username string    `json:"username"`
				FullName string    `json:"fullName"`
			}{
				ID:       user.ExternalID,
				Username: user.Username,
				FullName: user.FullName,
			}
		} else if message.SenderType == string(types.SenderTypeAgent) {
			var agent models.Agent
			initializers.DB.First(&agent, message.SenderID)
			sender = struct {
				ID     uuid.UUID `json:"id"`
				Name   string    `json:"name"`
				Lingo  string    `json:"lingo"`
				Traits []string  `json:"traits"`
			}{
				ID:     agent.ExternalID,
				Name:   agent.Name,
				Lingo:  agent.Lingo,
				Traits: agent.Traits,
			}
		}

		messagesWithSenders = append(messagesWithSenders, MessageWithSender{
			Message: message,
			Sender:  sender,
		})
	}

	// Prepare the chat response
	chatResponse := struct {
		models.Chat
		Messages []MessageWithSender `json:"messages"`
	}{
		Chat:     chat,
		Messages: messagesWithSenders,
	}

	c.JSON(http.StatusOK, gin.H{"data": chatResponse})
}

// CreateChatWithAgents creates a new chat with agents for the authenticated user
func CreateChatWithAgents(c *gin.Context) {
	var body struct {
		ChatName string `json:"chatName" binding:"required,max=20"`
		Agents   []struct {
			Name   string   `json:"name" binding:"required,max=20"`
			Lingo  string   `json:"lingo" binding:"required,max=20"`
			Traits []string `json:"traits" binding:"required"`
		} `json:"agents" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	chat := models.Chat{
		ChatName: body.ChatName,
		UserID:   userModel.ID,
	}

	result := initializers.DB.Create(&chat)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Create agents
	for _, agentData := range body.Agents {
		agent := models.Agent{
			Name:   agentData.Name,
			Lingo:  agentData.Lingo,
			Traits: agentData.Traits,
			ChatID: chat.ID,
		}

		if err := initializers.DB.Create(&agent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create agent"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": chat})
}
