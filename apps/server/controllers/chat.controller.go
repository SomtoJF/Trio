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

type createChatInput struct {
	ChatName string `json:"chatName" binding:"required,max=20"`
}

type addAgentToChatInput struct {
	Name   string   `json:"name" binding:"required,max=20"`
	Lingo  string   `json:"lingo" binding:"required,max=20"`
	Traits []string `json:"traits" binding:"required"`
}

type updateChatInput struct {
	ChatName string `json:"chatName" binding:"required,max=20"`
	Agents   []struct {
		ID     uuid.UUID `json:"id" binding:"required"`
		Name   string    `json:"name" binding:"required,max=20"`
		Lingo  string    `json:"lingo" binding:"required,max=20"`
		Traits []string  `json:"traits" binding:"required"`
	} `json:"agents" binding:"required"`
}

type createChatWithAgentsInput struct {
	ChatName string `json:"chatName" binding:"required,max=20"`
	Agents   []struct {
		Name   string   `json:"name" binding:"required,max=20"`
		Lingo  string   `json:"lingo" binding:"required,max=20"`
		Traits []string `json:"traits" binding:"required"`
	} `json:"agents" binding:"required"`
}

// CreateChat godoc
//
//	@Summary		Create a new chat
//	@Description	Creates a new chat for the authenticated user
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			chatName	body		createChatInput			true	"Chat name"
//	@Success		201			{object}	models.Chat				"Created chat"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		401			{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats [post]
func CreateChat(c *gin.Context) {
	var body createChatInput

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

// AddAgentToChat godoc
//
//	@Summary		Add an agent to a chat
//	@Description	Adds an agent to a chat (max 2 agents per chat)
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			chatId		path		string					true	"Chat ID"
//	@Param			agentInput	body		addAgentToChatInput		true	"Agent details"
//	@Success		201			{object}	models.Agent			"Created agent"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		401			{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404			{object}	map[string]interface{}	"Chat not found"
//	@Failure		409			{object}	map[string]interface{}	"Conflict"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId}/agents [post]
func AddAgentToChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body addAgentToChatInput

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

// DeleteChat godoc
//
//	@Summary		Delete a chat
//	@Description	Deletes a chat for the authenticated user
//	@Tags			chats
//	@Param			chatId	path		string					true	"Chat ID"
//	@Success		204		{object}	map[string]interface{}	"Chat deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		401		{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404		{object}	map[string]interface{}	"Chat not found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId} [delete]
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

// UpdateChat godoc
//
//	@Summary		Update a chat
//	@Description	Updates a chat's name and replaces agents with new ones
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			chatId		path		string					true	"Chat ID"
//	@Param			chatInput	body		updateChatInput			true	"Chat details"
//	@Success		200			{object}	models.Chat				"Updated chat"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		401			{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404			{object}	map[string]interface{}	"Chat not found"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId} [put]
func UpdateChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body updateChatInput

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

	// Update agents
	if err := initializers.DB.Model(&chat).Association("Agents").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear old agents"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"data": chat})
}

// GetChatInfo godoc
//
//	@Summary		Get chat information
//	@Description	Retrieves chat information including its agents and messages with sender details
//	@Tags			chats
//	@Param			chatId	path		string					true	"Chat ID"
//	@Success		200		{object}	models.Chat				"Chat information"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		401		{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404		{object}	map[string]interface{}	"Chat not found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId} [get]
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

// CreateChatWithAgents godoc
//
//	@Summary		Create a new chat with agents
//	@Description	Creates a new chat with agents for the authenticated user
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			chatInput	body		createChatWithAgentsInput	true	"Chat and agents details"
//	@Success		201			{object}	models.Chat					"Created chat with agents"
//	@Failure		400			{object}	map[string]interface{}		"Bad request"
//	@Failure		401			{object}	map[string]interface{}		"Unauthorized"
//	@Failure		500			{object}	map[string]interface{}		"Internal server error"
//	@Router			/chats/create-with-agents [post]
func CreateChatWithAgents(c *gin.Context) {
	var body createChatWithAgentsInput

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
