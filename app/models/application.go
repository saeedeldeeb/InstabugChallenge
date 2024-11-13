package models

import (
	"github.com/goravel/framework/database/orm"
)

type Application struct {
	orm.Model
	Token      string `gorm:"type:varchar(255);unique;not null"`
	Name       string `gorm:"type:varchar(255);not null"`
	ChatsCount int    `gorm:"type:int;default:0"`
}
