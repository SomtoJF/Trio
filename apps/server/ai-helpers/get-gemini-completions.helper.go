package aihelpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
)

func GetGeminiCompletions(c *gin.Context, request types.GeminiCompletionsRequest) (*genai.GenerateContentResponse, error) {
	client, ok := c.Value("GeminiClient").(*genai.Client)
	if !ok {
		return nil, fmt.Errorf("failed to get client from context")
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(c.Request.Context(), genai.Text(request.Prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response from Gemini")
	}

	inputTokens := int(resp.UsageMetadata.PromptTokenCount)
	outputTokens := int(resp.UsageMetadata.CandidatesTokenCount)
	totalTokens := inputTokens + outputTokens

	// Log the response to the database
	log := models.GeminiLogs{
		Prompt:       request.Prompt,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalTokens:  totalTokens,
		SenderType:   string(request.SenderType),
		SenderID:     request.SenderID,
	}

	result := initializers.DB.Create(&log)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to log Gemini response: %w", result.Error)
	}

	return resp, nil
}
