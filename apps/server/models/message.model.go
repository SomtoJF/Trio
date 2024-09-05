package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ExternalID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Content    string    `json:"content"`
	ChatID     uint      `json:"-"`
	SenderType string    `gorm:"type:sender_type_enum" json:"senderType"`
	SenderID   uint      `json:"_"`
	Chat       Chat      `gorm:"foreignKey:ChatID" json:"-"`
}
