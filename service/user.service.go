package service

import (
	"elearning_api/dto"
	"elearning_api/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type UserService interface {
	Confirm(request dto.UserConfirmRequest) error
	GetUserProfile(userID uint) (dto.UserBasicInfo, error)
}

type UserServiceImplement struct {
	UserRepository  repository.UserRepository
	TokenRepository repository.TokenRepository

	tx *gorm.DB
}

func NewUserService(userRepository repository.UserRepository, tokenRepository repository.TokenRepository) UserService {
	return &UserServiceImplement{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}

func (us UserServiceImplement) BeginTransaction(db *gorm.DB) *gorm.DB {
	us.tx = db.Begin()
	return us.tx
}

func (us UserServiceImplement) Confirm(request dto.UserConfirmRequest) error {
	userToken, err := us.TokenRepository.FindVerificationToken(request.Token)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiError(http.StatusBadRequest, "token not found")
		}

		return dto.NewApiError(http.StatusInternalServerError, "error when find token : "+err.Error())
	}

	isTokenExp := userToken.IsExpired()

	if isTokenExp {
		return dto.NewApiError(http.StatusBadRequest, "token is expired")
	}

	findUser, err := us.UserRepository.FindById(userToken.UserID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiError(http.StatusBadRequest, "user not found")
		}
		return dto.NewApiError(http.StatusInternalServerError, "error when find user by id : "+err.Error())
	}

	if findUser.PhoneNumberVerificationStatus {
		return dto.NewApiError(http.StatusBadRequest, "phone number already verified")
	}

	tx := us.BeginTransaction(us.UserRepository.GetDB())
	if tx.Error != nil {
		tx.Rollback()
		return dto.NewApiError(http.StatusInternalServerError, "error when start transaction : "+tx.Error.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("transaction panicked: %v", r)
		} else if tx.Error != nil {
			tx.Rollback()
			log.Fatalf("transaction failed: %v", tx.Error)
		}
	}()

	err = us.UserRepository.UpdatePhoneNumberVerificationStatus(dto.UserUpdateStatusVerifPhone{
		ID:                            userToken.UserID,
		PhoneNumberVerificationStatus: true,
	}, tx)

	if err != nil {
		tx.Rollback()
		return dto.NewApiError(http.StatusInternalServerError, "error when update phone number verification status : "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return dto.NewApiError(http.StatusInternalServerError, "error when commit transaction : "+err.Error())
	}

	return nil
}

func (us UserServiceImplement) GetUserProfile(userID uint) (dto.UserBasicInfo, error) {
	fmt.Println(userID)
	user, err := us.UserRepository.FindById(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserBasicInfo{}, dto.NewApiError(http.StatusBadRequest, "user not found")
		}

		return dto.UserBasicInfo{}, dto.NewApiError(http.StatusInternalServerError, "error when find user by id : "+err.Error())
	}

	return dto.UserBasicInfo{
		ID:                            user.ID,
		Username:                      user.Username,
		Email:                         user.Email,
		PhoneNumber:                   user.PhoneNumber,
		EmailVerificationStatus:       user.EmailVerificationStatus,
		PhoneNumberVerificationStatus: user.PhoneNumberVerificationStatus,
	}, nil
}
