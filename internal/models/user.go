package models

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Name              string
	CountryCode       string
	Image             string
	NotificationToken string
	Phone             string `gorm:"unique;index"`
}
