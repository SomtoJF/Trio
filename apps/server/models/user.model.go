package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ExternalID uuid.UUID `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"id"`
	Username   string    `gorm:"unique;type:string" json:"username"`
	FullName   string    `json:"fullName"`
	Chats      []Chat    `gorm:"foreignKey:UserID" json:"chats"`
}
