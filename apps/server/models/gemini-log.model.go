package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GeminiLogs struct {
	gorm.Model   `json:"-"`
	ExternalID   uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"id"`
	Prompt       string    `json:"prompt"`
	InputTokens  int       `json:"inputTokens"`
	OutputTokens int       `json:"outputTokens"`
	TotalTokens  int       `json:"totalTokens"`
	SenderType   string    `gorm:"type:sender_type_enum" json:"senderType"`
	SenderID     uint      `json:"_"`
}
