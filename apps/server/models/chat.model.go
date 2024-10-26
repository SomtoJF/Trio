package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatType string

const (
	ChatTypeDefault    ChatType = "DEFAULT"
	ChatTypeReflection ChatType = "REFLECTION"
)

type Chat struct {
	gorm.Model `json:"-"`
	ExternalID uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"id"`
	Messages   []Message `gorm:"foreignKey:ChatID" json:"messages"`
	UserID     uint      `json:"-"`
	ChatName   string    `json:"chatName"`
	Agents     []Agent   `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"agents"`
	Type       ChatType  `gorm:"type:varchar(11);check:type IN ('DEFAULT', 'REFLECTION');default:'DEFAULT'" json:"type"`
}
