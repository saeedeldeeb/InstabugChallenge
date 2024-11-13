package models

import (
	"github.com/goravel/framework/database/orm"
)

type Chat struct {
	orm.Model
	ApplicationId uint `gorm:"type:int;not null"`
	Application   *Application
	Number        int `gorm:"type:int;not null"`
	MessagesCount int `gorm:"type:int;default:0"`
}
