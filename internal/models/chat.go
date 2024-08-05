package models

import "gorm.io/gorm"

type Message struct {
	*gorm.Model
	Content     string
	AudioUrl    string
	ImageUrl    string
	UserID      uint `gorm:"index , foreignKey:UserID"`
	RecipientID uint `gorm:"index , foreignKey:RecipientID"` // Recipient user ID

}
