package repository

import (
	"elearning_api/model"
	"gorm.io/gorm"
)

type TokenRepository interface {
	CreateAccessToken(token model.AccessToken, tx *gorm.DB) error
	CreateAccessTokenLog(token model.AccessTokenLogs, tx *gorm.DB) error
	CreateRefreshToken(token model.RefreshToken, tx *gorm.DB) error
	CreateRefreshTokenLog(token model.RefreshTokenLogs, tx *gorm.DB) error
	CreateVerificationToken(token model.VerificationToken, tx *gorm.DB) error

	UpdateAccessToken(token model.AccessToken, tx *gorm.DB) error
	UpdateRefreshToken(token model.RefreshToken, tx *gorm.DB) error

	FindLastVerificationToken(userId uint) (model.VerificationToken, error)
	FindVerificationToken(token string) (model.VerificationToken, error)
	FindRefreshToken(token string) (model.RefreshToken, error)
	FindRefreshTokenByUserId(userId uint) (model.RefreshToken, error)
	FindRefreshTokenLog(token string) (model.RefreshTokenLogs, error)
	FindAccessToken(token string) (model.AccessToken, error)
	FindAccessTokenByUserId(userId uint) (model.AccessToken, error)
	FindAccessTokenLog(token string) (model.AccessTokenLogs, error)
}

type TokenRepositoryImplement struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &TokenRepositoryImplement{
		DB: db,
	}
}

func (tr TokenRepositoryImplement) CreateAccessToken(token model.AccessToken, tx *gorm.DB) error {
	return tx.Table("access_tokens").Create(&token).Error
}

func (tr TokenRepositoryImplement) CreateAccessTokenLog(token model.AccessTokenLogs, tx *gorm.DB) error {
	return tx.Table("access_token_logs").Create(&token).Error
}

func (tr TokenRepositoryImplement) CreateRefreshToken(token model.RefreshToken, tx *gorm.DB) error {
	return tx.Table("refresh_tokens").Create(&token).Error
}

func (tr TokenRepositoryImplement) CreateRefreshTokenLog(token model.RefreshTokenLogs, tx *gorm.DB) error {
	return tx.Table("refresh_token_logs").Create(&token).Error
}

func (tr TokenRepositoryImplement) CreateVerificationToken(token model.VerificationToken, tx *gorm.DB) error {
	return tx.Table("verification_tokens").Create(&token).Error
}

func (tr TokenRepositoryImplement) FindLastVerificationToken(userId uint) (model.VerificationToken, error) {
	var token model.VerificationToken
	err := tr.DB.Table("verification_tokens").Where("user_id = ?", userId).Order("created_at desc").First(&token).Error
	return token, err
}

func (tr TokenRepositoryImplement) FindRefreshToken(token string) (model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := tr.DB.Table("refresh_tokens").Where("token = ?", token).First(&refreshToken).Error
	return refreshToken, err
}

func (tr TokenRepositoryImplement) FindVerificationToken(token string) (model.VerificationToken, error) {
	var verificationToken model.VerificationToken
	err := tr.DB.Table("verification_tokens").Where("token = ?", token).First(&verificationToken).Error
	return verificationToken, err
}

func (tr TokenRepositoryImplement) UpdateAccessToken(token model.AccessToken, tx *gorm.DB) error {
	return tx.Table("access_tokens").Where("id = ?", token.ID).Updates(&token).Error
}

func (tr TokenRepositoryImplement) UpdateRefreshToken(token model.RefreshToken, tx *gorm.DB) error {
	return tx.Table("refresh_tokens").Where("id = ?", token.ID).Updates(&token).Error
}

func (tr TokenRepositoryImplement) FindRefreshTokenByUserId(userId uint) (model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := tr.DB.Table("refresh_tokens").Where("user_id = ?", userId).First(&refreshToken).Error
	return refreshToken, err
}

func (tr TokenRepositoryImplement) FindRefreshTokenLog(token string) (model.RefreshTokenLogs, error) {
	var refreshTokenLogs model.RefreshTokenLogs
	err := tr.DB.Table("refresh_token_logs").Where("token = ?", token).First(&refreshTokenLogs).Error
	return refreshTokenLogs, err
}

func (tr TokenRepositoryImplement) FindAccessToken(token string) (model.AccessToken, error) {
	var accessToken model.AccessToken
	err := tr.DB.Table("access_tokens").Where("token = ?", token).First(&accessToken).Error
	return accessToken, err
}

func (tr TokenRepositoryImplement) FindAccessTokenByUserId(userId uint) (model.AccessToken, error) {
	var accessToken model.AccessToken
	err := tr.DB.Table("access_tokens").Where("user_id = ?", userId).First(&accessToken).Error
	return accessToken, err
}

func (tr TokenRepositoryImplement) FindAccessTokenLog(token string) (model.AccessTokenLogs, error) {
	var accessTokenLogs model.AccessTokenLogs
	err := tr.DB.Table("access_token_logs").Where("token = ?", token).First(&accessTokenLogs).Error
	return accessTokenLogs, err
}
