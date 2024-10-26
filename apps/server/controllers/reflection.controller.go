package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
	"github.com/somtojf/trio/utils"
)

type reflectionMessageInput struct {
	Content string `json:"content" binding:"required"`
}

// PostReflectionMessage godoc
//
//	@Summary		Post a reflection message
//	@Description	Adds a reflection message to a chat for the authenticated user
//	@Tags			chat-messages
//	@Accept			json
//	@Produce		json
//	@Param			chatId			path		string					true	"Chat ID"
//	@Param			messageInput	body		reflectionMessageInput	true	"Message content"
//	@Success		200				{object}	map[string]interface{}	"Reflection message added successfully"
//	@Failure		400				{object}	map[string]interface{}	"Bad request"
//	@Failure		401				{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404				{object}	map[string]interface{}	"Chat not found"
//	@Failure		500				{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId}/messages/reflection [post]
func PostReflectionMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body reflectionMessageInput

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

	client, ok := c.Value("GeminiClient").(*genai.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't retrieve gemini client"})
		return
	}

	userMessage := models.Message{
		Content:    body.Content,
		SenderType: string(types.SenderTypeUser),
		SenderID:   userModel.ID,
		ChatID:     chat.ID,
	}

	if err := initializers.DB.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user message to chat"})
		return
	}

	chatHistory, err := utils.GetChatHistory(chat.ID, utils.MAX_TOKENS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat history"})
		return
	}

	shuffledAgents := utils.RandomizeArrayElements(chat.Agents)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	responseChan := make(chan utils.AgentResponse, len(shuffledAgents))
	doneChan := make(chan struct{})

	// Start the agent response loop
	go utils.AgentResponseLoop(c.Request.Context(), client, shuffledAgents, chatHistory, body.Content, responseChan, doneChan)

	// Stream responses to the client
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				return
			}
			data, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				continue
			}
			c.SSEvent("message", string(data))
			c.Writer.Flush()
		case <-doneChan:
			return
		}
	}
}
