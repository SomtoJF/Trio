package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model `json:"-"`
	ExternalID uuid.UUID      `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string         `json:"name"`
	ChatID     uint           `json:"-"`
	Lingo      string         `json:"lingo"`
	Traits     pq.StringArray `gorm:"type:text[]" json:"traits"`
}
