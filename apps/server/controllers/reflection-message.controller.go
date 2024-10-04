package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
	"github.com/somtojf/trio/utils"
)

func PostReflectionMessage(c *gin.Context) {
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

	responseChan := make(chan AgentResponse, len(shuffledAgents))
	doneChan := make(chan struct{})

	// Start the agent response loop
	go agentResponseLoop(c.Request.Context(), client, shuffledAgents, chatHistory, body.Content, responseChan, doneChan)

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

func agentResponseLoop(ctx context.Context, client *genai.Client, agents []models.Agent, chatHistory []models.Message, userMessage string, responseChan chan<- AgentResponse, doneChan chan<- struct{}) {
	defer close(responseChan)
	defer close(doneChan)

	agentResponses := make(map[uint]string)
	for {
		for _, agent := range agents {
			response := generateAgentResponseAsync(ctx, client, agent, chatHistory, userMessage, agentResponses)
			agentResponses[agent.ID] = response

			responseChan <- AgentResponse{
				AgentID:   agent.ID,
				AgentName: agent.Name,
				Content:   response,
			}

			if strings.HasPrefix(strings.ToLower(response), "agree") || strings.HasSuffix(strings.ToLower(response), "alternate") || response == "" {
				return
			}
		}

		for agentID, response := range agentResponses {
			chatHistory = append(chatHistory, models.Message{
				Content:    response,
				SenderType: string(types.SenderTypeAgent),
				SenderID:   agentID,
				ChatID:     chatHistory[0].ChatID,
			})
		}
	}
}

func generateAgentResponseAsync(ctx context.Context, client *genai.Client, agent models.Agent, chatHistory []models.Message, userMessage string, otherAgentResponses map[uint]string) string {
	model := client.GenerativeModel("gemini-1.5-pro")

	prompt := generateReflectionPrompt(agent, chatHistory, userMessage, otherAgentResponses)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error generating content for agent %s: %v", agent.Name, err)
		return ""
	}

	response := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(response)
}

func generateReflectionPrompt(agent models.Agent, chatHistory []models.Message, userMessage string, otherAgentResponses map[uint]string) string {
	prompt := fmt.Sprintf(`
	You are %s, a helpful AI agent with freedom to provide responses in the best way you see fit.
	You are in a group chat with a human user and another AI agent. Your goal is to collaborate with the other agent to respond to the user's message. You're response has three states depending on your alignment with the other agent's response:
	1. You can agree with the other agent's response, in which case you must respond ONLY with the word "agree".
	2. You can disagree with the other agent's message, in which case you present your opposing view.
	3. You can contribute to the other agent's message, which means that partially agree with their response and present your own view. You must also end your response with the word "alternate".

	If there has been no response to the user's message, you can respond with a solution which you deem fit.
	Chat History:
	%s

	The user's latest message is: "%s"
	`, agent.Name, utils.FormatChatHistory(chatHistory), userMessage)

	// Add information about other agents' responses to the prompt
	for agentID, response := range otherAgentResponses {
		if agentID != agent.ID {
			prompt += fmt.Sprintf("\nThe other agent's response: %s", response)
		}
	}

	return prompt
}

type AgentResponse struct {
	AgentID   uint   `json:"agentId"`
	AgentName string `json:"agentName"`
	Content   string `json:"content"`
}
