package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
)

func AgentResponseLoop(ctx context.Context, client *genai.Client, agents []models.Agent, chatHistory []models.Message, userMessage string, responseChan chan<- AgentResponse, doneChan chan<- struct{}) {
	defer close(responseChan)
	defer close(doneChan)

	agentResponses := make(map[uint]string)
	for {
		for _, agent := range agents {
			response := GenerateAgentResponseAsync(ctx, client, agent, chatHistory, userMessage, agentResponses)
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

func GenerateAgentResponseAsync(ctx context.Context, client *genai.Client, agent models.Agent, chatHistory []models.Message, userMessage string, otherAgentResponses map[uint]string) string {
	model := client.GenerativeModel("gemini-1.5-pro")

	prompt := GenerateReflectionPrompt(agent, chatHistory, userMessage, otherAgentResponses)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error generating content for agent %s: %v", agent.Name, err)
		return ""
	}

	response := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(response)
}

func GenerateReflectionPrompt(agent models.Agent, chatHistory []models.Message, userMessage string, otherAgentResponses map[uint]string) string {
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
	`, agent.Name, FormatChatHistory(chatHistory), userMessage)

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
