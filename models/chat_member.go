package models

import (
	_ "encoding/gob"
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ChatMember struct {
	DateCreated  time.Time
	DateModified time.Time
	ID           int64  `json:"id"`
	ChatID       int64  `json:"chat_id"`
	UserID       int64  `json:"user_id"`
	UserName     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

func (ChatMember) TableName() string {
	return "chat_member"
}
