package types

type GeminiCompletionsRequest struct {
	Prompt     string
	SenderID   uint
	SenderType SenderType
}
