package clients

import (
	"context"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func CreateGeminiClient() (*genai.Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client, nil
}
