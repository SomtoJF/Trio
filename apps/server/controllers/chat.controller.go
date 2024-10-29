package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/response"
	"github.com/somtojf/trio/types"
	"github.com/somtojf/trio/utils"
	"gorm.io/gorm"
)

type addAgentToChatInput struct {
	Name   string   `json:"name" binding:"required,max=20"`
	Lingo  string   `json:"lingo" binding:"required,max=20"`
	Traits []string `json:"traits" binding:"required"`
}

type addMessageToChatInput struct {
	Content string `json:"content" binding:"required"`
}

type updateChatInput struct {
	ChatName string `json:"chatName" binding:"required,max=20"`
	Agents   []struct {
		ID       uuid.UUID `json:"id" binding:"required"`
		Name     string    `json:"name" binding:"required,max=20"`
		Metadata struct {
			Lingo  string   `json:"lingo" binding:"required,max=20"`
			Traits []string `json:"traits" binding:"required"`
		}
	} `json:"agents" binding:"required"`
}

type createChatWithAgentsInput struct {
	ChatName string `json:"chatName" binding:"required,max=20"`
	Type     string `json:"type" binding:"oneof=DEFAULT REFLECTION"`
	Agents   []struct {
		Name   string   `json:"name" binding:"required,max=20"`
		Lingo  string   `json:"lingo" binding:"required,max=20"`
		Traits []string `json:"traits" binding:"required"`
	} `json:"agents" binding:"required"`
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

	var agentMetadata *models.AgentMetadata

	if chat.Type == models.ChatTypeReflection {
		agentMetadata = nil
	} else {
		agentMetadata = &models.AgentMetadata{
			Lingo:  body.Lingo,
			Traits: body.Traits,
		}
	}

	agent := models.Agent{
		Name:     body.Name,
		Metadata: agentMetadata,
		ChatID:   chat.ID,
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
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

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

	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	chat.ChatName = body.ChatName

	if err := tx.Save(&chat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	if err := tx.Model(&chat).Association("Agents").Clear(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear old agents"})
		return
	}

	var agentMetadata []*models.AgentMetadata

	if chat.Type == models.ChatTypeReflection {
		if len(body.Agents) != 2 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reflection chat must have exactly two agents"})
			return
		}
		for range body.Agents {
			agentMetadata = append(agentMetadata, nil)
		}
	} else {
		for _, agentData := range body.Agents {
			agentMetadata = append(agentMetadata, &models.AgentMetadata{
				Lingo:  agentData.Metadata.Lingo,
				Traits: agentData.Metadata.Traits,
			})
		}
	}

	for i, agentData := range body.Agents {
		agent := models.Agent{
			Name:     agentData.Name,
			Metadata: agentMetadata[i],
			ChatID:   chat.ID,
		}

		if err := tx.Create(&agent).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create agent"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
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
			initializers.DB.Preload("Metadata").First(&agent, message.SenderID)
			sender = struct {
				ID       uuid.UUID `json:"id"`
				Name     string    `json:"name"`
				Metadata struct {
					Lingo  string   `json:"lingo"`
					Traits []string `json:"traits"`
				}
			}{
				ID:   agent.ExternalID,
				Name: agent.Name,
				Metadata: struct {
					Lingo  string   `json:"lingo"`
					Traits []string `json:"traits"`
				}{
					Lingo:  agent.Metadata.Lingo,
					Traits: agent.Metadata.Traits,
				},
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
func CreateChat(c *gin.Context) {
	var body createChatWithAgentsInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, agent := range body.Agents {
		if len(agent.Traits) > 4 || len(agent.Traits) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Agent must have at least one trait and a maximum of four traits"})
			return
		}
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	tx := initializers.DB.Begin()

	chat := models.Chat{
		ChatName: body.ChatName,
		Type:     models.ChatType(body.Type),
		UserID:   userModel.ID,
	}

	if err := tx.Create(&chat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	if body.Type == string(models.ChatTypeReflection) {
		if len(body.Agents) != 2 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reflection chat must have exactly two agents"})
			return
		}

		for _, agent := range body.Agents {
			if err := createAgent(tx, agent.Name, nil, chat.ID); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		if len(body.Agents) == 0 || len(body.Agents) > 2 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Chat must have at least one agent and a maximum of two agents"})
			return
		}

		for _, agent := range body.Agents {
			agentMetadata := &models.AgentMetadata{
				Lingo:  agent.Lingo,
				Traits: agent.Traits,
			}
			if err := createAgent(tx, agent.Name, agentMetadata, chat.ID); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": chat})
}

// Helper function to create an agent
func createAgent(tx *gorm.DB, name string, metadata *models.AgentMetadata, chatID uint) error {
	agent := models.Agent{
		Name:     name,
		Metadata: metadata,
		ChatID:   chatID,
	}
	return tx.Create(&agent).Error
}

// NewMessage godoc
//
//	@Summary		Add a new message to a chat
//	@Description	Adds a new message to a chat and generates responses from agents
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			chatId		path		string					true	"Chat ID"
//	@Param			messageInput	body	addMessageToChatInput	true	"Message content"
//	@Success		201			{object}	map[string]interface{}	"Message added successfully"
//	@Success		200			{object}	map[string]interface{}	"Reflection response generated successfully"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		401			{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404			{object}	map[string]interface{}	"Chat not found"
//	@Failure		424			{object}	map[string]interface{}	"Chat must have at least one agent"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId}/messages [post]
func NewMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body addMessageToChatInput

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
	if err := initializers.DB.Preload("Agents").First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	if len(chat.Agents) == 0 {
		c.JSON(http.StatusFailedDependency, gin.H{"error": "Chat must have at least one agent"})
		return
	}

	client, ok := c.Value("GeminiClient").(*genai.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't retrieve gemini client"})
		return
	}

	// Ensure we have at least one agent
	if len(chat.Agents) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat must have at least one agent"})
		return
	}

	response := response.NewResponse(chat.Messages, chat, chat.Agents, userModel, c, client)
	if chat.Type == models.ChatTypeDefault {
		agentResponses, err := response.GenerateBasicResponse(body.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := utils.SaveResponsesToDatabase(agentResponses...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save agent responses"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"requestPrompt": body.Content,
			"data":          agentResponses,
		})
		return
	} else if chat.Type == models.ChatTypeReflection {
		err = response.GenerateReflectionResponse(body.Content)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DeleteAllChats godoc
//
//	@Summary		Delete all chats for the authenticated user
//	@Description	Deletes all chats and associated data belonging to the authenticated user
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"All chats deleted successfully"
//	@Failure		401	{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats [delete]
func DeleteAllChats(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	tx := initializers.DB.Begin()

	if err := tx.Where("user_id = ?", userModel.ID).Delete(&models.Chat{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chats"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All chats deleted successfully"})
}
