package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model `json:"-"`
	ExternalID uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"id"`
	Messages   []Message `gorm:"foreignKey:ChatID" json:"messages"`
	UserID     uint      `json:"-"`
	ChatName   string    `json:"chatName"`
	Agents     []Agent   `gorm:"foreignKey:ChatID" json:"agents"`
}
