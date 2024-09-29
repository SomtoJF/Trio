package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
)

const (
	MAX_TOKENS = 4000
)

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
	if err := initializers.DB.Preload("Agents").First(&chat, "external_id = ? AND user_id = ?", chatID, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	if len(chat.Agents) != 2 {
		c.JSON(http.StatusFailedDependency, gin.H{"error": "Chat must have exactly two agents"})
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

	// Get chat history
	chatHistory, err := getChatHistory(chat.ID, MAX_TOKENS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat history"})
		return
	}

	// Randomize agent response order
	firstAgent, secondAgent := randomizeAgents(chat.Agents)

	// Generate responses
	firstResponse, err := generateAgentResponse(c.Request.Context(), client, firstAgent, chatHistory, body.Content, userModel.Username, secondAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate response for %s", firstAgent.Name)})
		return
	}

	// Add first agent's response to chat history
	chatHistory = append(chatHistory, firstResponse)

	secondResponse, err := generateAgentResponse(c.Request.Context(), client, secondAgent, chatHistory, body.Content, userModel.Username, firstAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate response for %s", secondAgent.Name)})
		return
	}

	// Save responses to database
	if err := saveResponsesToDatabase(firstResponse, secondResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save agent responses"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"requestPrompt": userMessage,
		"data":          []models.Message{firstResponse, secondResponse},
	})
}

func getChatHistory(chatID uint, maxTokens int) ([]models.Message, error) {
	var messages []models.Message
	err := initializers.DB.Where("chat_id = ?", chatID).Order("created_at DESC").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// Reverse the order to get chronological order
	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}

	// TODO: Implement token counting and truncation logic here
	// For now, we'll just return all messages
	return messages, nil
}

func randomizeAgents(agents []models.Agent) (models.Agent, models.Agent) {
	if rand.Float32() < 0.5 {
		return agents[0], agents[1]
	}
	return agents[1], agents[0]
}

func generateAgentResponse(ctx context.Context, client *genai.Client, agent models.Agent, chatHistory []models.Message, userMessage string, userName string, otherAgent models.Agent) (models.Message, error) {
	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := createEnhancedPrompt(agent, chatHistory, userMessage, userName, otherAgent)

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return models.Message{}, err
	}

	var aiResponse string
	if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
		aiResponse = string(res.Candidates[0].Content.Parts[0].(genai.Text))
	} else {
		aiResponse = "No response generated"
	}

	return models.Message{
		Content:    aiResponse,
		SenderType: string(types.SenderTypeAgent),
		SenderID:   agent.ID,
		ChatID:     chatHistory[0].ChatID,
	}, nil
}

func createEnhancedPrompt(agent models.Agent, chatHistory []models.Message, userMessage string, userName string, otherAgent models.Agent) string {
	return fmt.Sprintf(`
You are %s, an AI agent with the following traits: %s.
You are in a group chat with a human user called %s and another AI agent named %s with traits: %s.
Chat History:
%s

The user's latest message is: "%s" 

Please respond to the user's message and, if appropriate, to the other agent's previous message. Refer to them as @<targetname>.
Use your defined traits to guide your response style and content.
Engage in a natural, flowing conversation while keeping responses as short as possible, and feel free to ask questions or make observations to keep the dialogue engaging.
Remember as much context as you can from previous messages and use them when necessary.
`, agent.Name, strings.Join(agent.Traits, ", "), userName, otherAgent.Name, strings.Join(otherAgent.Traits, ", "),
		formatChatHistory(chatHistory), userMessage)
}

func formatChatHistory(history []models.Message) string {
	var formattedHistory strings.Builder
	for _, msg := range history {
		formattedHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.SenderType, msg.Content))
	}
	return formattedHistory.String()
}

func saveResponsesToDatabase(responses ...models.Message) error {
	return initializers.DB.Create(&responses).Error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
