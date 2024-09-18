package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
	"gorm.io/gorm"
)

// AddMessageToChat adds a new message to a chat
func AddMessageToChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var body struct {
		Content string `json:"content" binding:"required"`
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
	if err := initializers.DB.Preload("Agents").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(10)
	}).First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	if len(chat.Agents) < 1 {
		c.JSON(http.StatusFailedDependency, gin.H{"error": "No agents found in the chat"})
		return
	}

	client, ok := c.Value("GeminiClient").(*genai.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't retrieve gemini client"})
		return
	}

	// Create and save the user's message
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

	var wg sync.WaitGroup
	responses := make([]models.Message, len(chat.Agents))

	for i, agent := range chat.Agents {
		wg.Add(1)
		go func(i int, agent models.Agent) {
			defer wg.Done()

			model := client.GenerativeModel("gemini-1.5-flash")
			cs := model.StartChat()

			// Populate cs.History with the last 10 messages (in chronological order)
			for j := len(chat.Messages) - 1; j >= 0; j-- {
				msg := chat.Messages[j]
				role := "user"
				if msg.SenderType == string(types.SenderTypeAgent) {
					role = "model"
				}
				cs.History = append(cs.History, &genai.Content{
					Parts: []genai.Part{
						genai.Text(msg.Content),
					},
					Role: role,
				})
			}

			res, err := cs.SendMessage(c.Request.Context(), genai.Text(body.Content))
			if err != nil {
				log.Printf("Error getting response from agent %s: %v", agent.Name, err)
				return
			}

			var aiResponse string
			if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
				aiResponse = string(res.Candidates[0].Content.Parts[0].(genai.Text))
			} else {
				aiResponse = "No response generated"
			}

			responses[i] = models.Message{
				Content:    aiResponse,
				SenderType: string(types.SenderTypeAgent),
				SenderID:   agent.ID,
				ChatID:     chat.ID,
			}
		}(i, agent)
	}

	wg.Wait()

	// Save AI responses to the database
	for _, response := range responses {
		if err := initializers.DB.Create(&response).Error; err != nil {
			log.Printf("Error saving AI response to database: %v", err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"requestPrompt": userMessage,
		"data":          responses,
	})
}
