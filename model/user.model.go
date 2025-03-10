package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID                            uint               `gorm:"primaryKey"`
	Username                      string             `gorm:"unique;not null"`
	Password                      string             `gorm:"not null"`
	Email                         string             `gorm:"not null"`
	EmailVerificationStatus       bool               `gorm:"default:false"`
	PhoneNumber                   string             `gorm:"not null;default:''"`
	PhoneNumberVerificationStatus bool               `gorm:"default:false"`
	RegistrationMethodID          uint               `gorm:"not null"`
	RegistrationMethod            RegistrationMethod `gorm:"foreignKey:RegistrationMethodID"`

	AccessToken  AccessToken  `gorm:"foreignKey:UserID"`
	RefreshToken RefreshToken `gorm:"foreignKey:UserID"`
}

type RegistrationMethod struct {
	gorm.Model

	ID     uint   `gorm:"primaryKey"`
	Method string `gorm:"unique;not null"`
	Users  []User `gorm:"foreignKey:RegistrationMethodID"`
}
