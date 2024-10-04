package utils

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
)

const (
	MAX_TOKENS = 4000
)

func RandomizeArrayElements[T any](array []T) []T {
	shuffled := make([]T, len(array))
	copy(shuffled, array)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func GetChatHistory(chatID uint, maxTokens int) ([]models.Message, error) {
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

func FormatChatHistory(history []models.Message) string {
	var formattedHistory strings.Builder
	for _, msg := range history {
		formattedHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.SenderType, msg.Content))
	}
	return formattedHistory.String()
}
