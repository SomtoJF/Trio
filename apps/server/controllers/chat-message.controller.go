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
	"github.com/somtojf/trio/utils"
)

type addMessageToChatInput struct {
	Content string `json:"content" binding:"required"`
}

// AddMessageToChat godoc
//
//	@Summary		Add a message to a chat
//	@Description	Adds a message to a chat for the authenticated user
//	@Tags			chat-messages
//	@Accept			json
//	@Produce		json
//	@Param			chatId			path		string					true	"Chat ID"
//	@Param			messageInput	body		addMessageToChatInput	true	"Message content"
//	@Success		201				{object}	map[string]interface{}	"Message added successfully"
//	@Failure		400				{object}	map[string]interface{}	"Bad request"
//	@Failure		401				{object}	map[string]interface{}	"Unauthorized"
//	@Failure		404				{object}	map[string]interface{}	"Chat not found"
//	@Failure		500				{object}	map[string]interface{}	"Internal server error"
//	@Router			/chats/{chatId}/messages [post]
func AddMessageToChat(c *gin.Context) {
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
	chatHistory, err := utils.GetChatHistory(chat.ID, utils.MAX_TOKENS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat history"})
		return
	}

	shuffledAgents := utils.RandomizeArrayElements(chat.Agents)

	var agentResponses []models.Message

	// Generate responses
	for i, agent := range shuffledAgents {
		var otherAgent models.Agent
		if i+1 < len(shuffledAgents) {
			otherAgent = shuffledAgents[i+1]
		} else if len(shuffledAgents) > 1 {
			otherAgent = shuffledAgents[0]
		}

		response, err := generateAgentResponse(c.Request.Context(), client, agent, chatHistory, body.Content, userModel.Username, otherAgent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate response for %s", agent.Name)})
			return
		}

		agentResponses = append(agentResponses, response)
		chatHistory = append(chatHistory, response)
	}

	// Save responses to database
	if err := saveResponsesToDatabase(agentResponses...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save agent responses"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"requestPrompt": userMessage,
		"data":          agentResponses,
	})
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
`, agent.Name, strings.Join(agent.Metadata.Traits, ", "), userName, otherAgent.Name, strings.Join(otherAgent.Metadata.Traits, ", "),
		utils.FormatChatHistory(chatHistory), userMessage)
}

func saveResponsesToDatabase(responses ...models.Message) error {
	return initializers.DB.Create(&responses).Error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
