package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model `json:"-"`
	ExternalID uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string    `json:"name"`
	ChatID     uint      `json:"-"`
	Lingo      string    `json:"lingo"`
	Traits     []string  `gorm:"type:text[]" json:"traits"`
}
