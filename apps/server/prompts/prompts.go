package prompts

import (
	"fmt"
	"strings"

	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/utils"
)

type Prompt struct {
	agent       models.Agent
	chatHistory []models.Message
	userName    string
	otherAgent  models.Agent
	userMessage string
}

func NewPromptGenerator(agent models.Agent, chatHistory []models.Message, userName string, otherAgent models.Agent, userMessage string) Prompt {
	return Prompt{
		agent:       agent,
		chatHistory: chatHistory,
		userName:    userName,
		otherAgent:  otherAgent,
		userMessage: userMessage,
	}
}

func (p *Prompt) GenerateBasicPrompt() string {
	var otherAgentTraits string = ""
	if p.otherAgent.Metadata != nil {
		otherAgentTraits = strings.Join(p.otherAgent.Metadata.Traits, ", ")
	}
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
`, p.agent.Name, strings.Join(p.agent.Metadata.Traits, ", "), p.userName, p.otherAgent.Name, otherAgentTraits,
		utils.FormatChatHistory(p.chatHistory), p.userMessage)
}

func (p *Prompt) GenerateReflectionPrompt(otherAgentResponses map[uint]string) string {
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
	`, p.agent.Name, utils.FormatChatHistory(p.chatHistory), p.userMessage)

	// Add information about other agents' responses to the prompt
	for agentID, response := range otherAgentResponses {
		if agentID != p.agent.ID {
			prompt += fmt.Sprintf("\nThe other agent's response: %s", response)
		}
	}

	return prompt
}
