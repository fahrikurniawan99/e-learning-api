package model

import (
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	gorm.Model

	ID    uint   `gorm:"primaryKey"`
	Token string `gorm:"not null"`

	UserID uint `gorm:"not null;unique"`
}

type RefreshToken struct {
	gorm.Model

	ID    uint   `gorm:"primaryKey"`
	Token string `gorm:"not null"`

	UserID uint `gorm:"not null;unique"`
}

type VerificationToken struct {
	gorm.Model

	ID    uint   `gorm:"primaryKey"`
	Token string `gorm:"not null"`

	UserID uint `gorm:"not null"`
}

func (verif VerificationToken) IsExpired() bool {
	if verif.CreatedAt.IsZero() {
		return true
	}

	duration := time.Since(verif.CreatedAt)

	return duration > 5*time.Minute
}

type AccessTokenLogs struct {
	gorm.Model

	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null"`
}

type RefreshTokenLogs struct {
	gorm.Model

	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null"`
}
