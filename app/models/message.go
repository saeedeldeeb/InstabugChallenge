package models

import (
	"github.com/goravel/framework/database/orm"
)

type Message struct {
	orm.Model
	ChatId uint `gorm:"type:int;not null"`
	Chat   *Chat
	Number int    `gorm:"type:int;not null"`
	Body   string `gorm:"type:text;not null"`
}
