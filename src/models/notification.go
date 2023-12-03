// models/notification.go
package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	SenderId 	uint 
	Sender 		User `gorm:"foreignKey:SenderId"`
	ReceiverId 	uint
	Receiver 	User `gorm:"foreignKey:ReceiverId"`
	Content 	string
}