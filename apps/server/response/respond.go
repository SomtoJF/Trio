package response

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/prompts"
	"github.com/somtojf/trio/types"
	"github.com/somtojf/trio/utils"
)

type ReflectionAgentResponse struct {
	AgentID   uint   `json:"agentId"`
	AgentName string `json:"agentName"`
	Content   string `json:"content"`
}

type Response struct {
	ChatHistory []models.Message
	Chat        models.Chat
	Agents      []models.Agent
	User        models.User
	Context     *gin.Context
	Client      *genai.Client
}

func NewResponse(chatHistory []models.Message, chat models.Chat, agents []models.Agent, user models.User, context *gin.Context, client *genai.Client) Response {
	return Response{
		ChatHistory: chatHistory,
		Chat:        chat,
		Agents:      agents,
		User:        user,
		Context:     context,
		Client:      client,
	}
}

func (r *Response) GenerateBasicResponse(prompt string) ([]models.Message, error) {
	userMessage := models.Message{
		Content:    prompt,
		SenderType: string(types.SenderTypeUser),
		SenderID:   r.User.ID,
		ChatID:     r.Chat.ID,
	}

	if err := initializers.DB.Create(&userMessage).Error; err != nil {
		return nil, fmt.Errorf("Failed to add user message to chat")
	}

	// Get chat history
	chatHistory, err := utils.GetChatHistory(r.Chat.ID, utils.MAX_TOKENS)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve chat history")
	}

	shuffledAgents := utils.RandomizeArrayElements(r.Chat.Agents)

	var agentResponses []models.Message

	// Generate responses
	for i, agent := range shuffledAgents {
		var otherAgent models.Agent
		if i+1 < len(shuffledAgents) {
			otherAgent = shuffledAgents[i+1]
		} else if len(shuffledAgents) > 1 {
			otherAgent = shuffledAgents[0]
		}

		response, err := generateAgentResponse(r.Context.Request.Context(), r.Client, agent, chatHistory, prompt, r.User.Username, otherAgent)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate response for %s", agent.Name)
		}

		agentResponses = append(agentResponses, response)
		chatHistory = append(chatHistory, response)
	}

	return agentResponses, nil
}

func (r *Response) GenerateReflectionResponse(prompt string) error {
	userMessage := models.Message{
		Content:    prompt,
		SenderType: string(types.SenderTypeUser),
		SenderID:   r.User.ID,
		ChatID:     r.Chat.ID,
	}

	if err := initializers.DB.Create(&userMessage).Error; err != nil {
		return fmt.Errorf("Failed to add user message to chat")
	}

	chatHistory, err := utils.GetChatHistory(r.Chat.ID, utils.MAX_TOKENS)
	if err != nil {
		return fmt.Errorf("Failed to retrieve chat history")
	}

	shuffledAgents := utils.RandomizeArrayElements(r.Chat.Agents)

	r.Context.Writer.Header().Set("Content-Type", "text/event-stream")
	r.Context.Writer.Header().Set("Cache-Control", "no-cache")
	r.Context.Writer.Header().Set("Connection", "keep-alive")
	r.Context.Writer.Header().Set("Transfer-Encoding", "chunked")

	responseChan := make(chan ReflectionAgentResponse, len(shuffledAgents))
	doneChan := make(chan struct{})

	// Start the agent response loop
	go ReflectionAgentResponseLoop(r.Context.Request.Context(), r.Client, shuffledAgents, chatHistory, userMessage.Content, responseChan, doneChan)

	// Stream responses to the client
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				return nil
			}
			data, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				continue
			}
			r.Context.SSEvent("message", string(data))
			r.Context.Writer.Flush()
		case <-doneChan:
			return nil
		}
	}
}

func ReflectionAgentResponseLoop(ctx context.Context, client *genai.Client, agents []models.Agent, chatHistory []models.Message, userMessage string, responseChan chan<- ReflectionAgentResponse, doneChan chan<- struct{}) {
	defer close(responseChan)
	defer close(doneChan)

	agentResponses := make(map[uint]string)
	for {
		for _, agent := range agents {
			response := GenerateAgentResponseAsync(ctx, client, agent, chatHistory, userMessage, agentResponses)
			agentResponses[agent.ID] = response

			responseChan <- ReflectionAgentResponse{
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

func GenerateAgentResponseAsync(ctx context.Context, client *genai.Client, agent models.Agent, chatHistory []models.Message, userMessage string, otherAgentResponses map[uint]string) string {
	model := client.GenerativeModel("gemini-1.5-pro")
	promptGenerator := prompts.NewPromptGenerator(agent, chatHistory, "", models.Agent{}, userMessage)

	prompt := promptGenerator.GenerateReflectionPrompt(otherAgentResponses)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error generating content for agent %s: %v", agent.Name, err)
		return ""
	}

	response := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(response)
}

func generateAgentResponse(ctx context.Context, client *genai.Client, agent models.Agent, chatHistory []models.Message, userMessage string, userName string, otherAgent models.Agent) (models.Message, error) {
	model := client.GenerativeModel("gemini-1.5-flash")
	promptGenerator := prompts.NewPromptGenerator(agent, chatHistory, userName, otherAgent, userMessage)
	prompt := promptGenerator.GenerateBasicPrompt()

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

func init() {
	rand.Seed(time.Now().UnixNano())
}
