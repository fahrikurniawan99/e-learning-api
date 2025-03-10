package repository

import (
	"elearning_api/dto"
	"elearning_api/model"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetDB() *gorm.DB
	FindByEmailAddress(email string) (model.User, error)
	Create(user *model.User, tx *gorm.DB) error
	FindByPhoneNumber(phoneNumber string) (model.User, error)
	UpdatePhoneNumberVerificationStatus(user dto.UserUpdateStatusVerifPhone, tx *gorm.DB) error
	FindById(id uint) (model.User, error)
}

type UserRepositoryImplement struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImplement{
		DB: db,
	}
}

func (ur UserRepositoryImplement) GetDB() *gorm.DB {
	return ur.DB
}

func (ur UserRepositoryImplement) FindByEmailAddress(email string) (model.User, error) {
	fmt.Println(email)
	var user model.User
	err := ur.DB.Table("users").Where("email = ?", email).Select("id, username, email, phone_number, phone_number_verification_status, email_verification_status").First(&user).Error
	return user, err
}

func (ur UserRepositoryImplement) Create(user *model.User, tx *gorm.DB) error {
	return tx.Table("users").Create(&user).Error
}

func (ur UserRepositoryImplement) FindByPhoneNumber(phoneNumber string) (model.User, error) {
	var user model.User
	err := ur.DB.Table("users").Where("phone_number = ?", phoneNumber).Select("id, username, email, phone_number, phone_number_verification_status, email_verification_status, password").First(&user).Error
	return user, err
}

func (ur UserRepositoryImplement) UpdatePhoneNumberVerificationStatus(user dto.UserUpdateStatusVerifPhone, tx *gorm.DB) error {
	return tx.Table("users").Where("id = ?", user.ID).Update("phone_number_verification_status", user.PhoneNumberVerificationStatus).Error
}

func (ur UserRepositoryImplement) FindById(id uint) (model.User, error) {
	var user model.User
	err := ur.DB.Table("users").Where("id = ?", id).Select("id, username, email, phone_number, phone_number_verification_status, email_verification_status").First(&user).Error
	return user, err
}
