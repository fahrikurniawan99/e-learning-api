package service

import (
	"elearning_api/dto"
	"elearning_api/model"
	"elearning_api/repository"
	"elearning_api/utils"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type AuthService interface {
	SocialLogin(request dto.SocialLoginRequest) (dto.JwtToken, error)
	ManualRegister(request dto.ManualRegisterRequest) error
	ManualLogin(request dto.ManualLoginRequest) (dto.JwtToken, error)
	SendVerificationLink(request dto.SendVerificationLinkRequest) error
	RefreshToken(request dto.RefreshTokenRequest) (dto.JwtToken, error)
	createAccessTokens(userID uint) (string, error)
	createRefreshToken(userID uint) (string, error)
}

type AuthServiceImplement struct {
	UserRepository  repository.UserRepository
	TokenRepository repository.TokenRepository

	GoogleService GoogleService
	WaService     WhatsappService
	tx            *gorm.DB
}

func NewAuthService(userRepository repository.UserRepository, tokenRepository repository.TokenRepository, googleService GoogleService, wa WhatsappService) AuthService {
	return &AuthServiceImplement{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
		GoogleService:   googleService,
		WaService:       wa,
	}
}

func (as *AuthServiceImplement) BeginTransaction(db *gorm.DB) *gorm.DB {
	as.tx = db.Begin()
	return as.tx
}

func (as AuthServiceImplement) SocialLogin(request dto.SocialLoginRequest) (dto.JwtToken, error) {
	googleToken, err := as.GoogleService.GoogleVerifyToken(request.Code)

	if err != nil {
		return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Error when verify google token : "+err.Error())
	}

	userInfo, err := as.GoogleService.GetGoogleUserInfo(googleToken.IDToken)
	if err != nil {
		return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Error when get user info : "+err.Error())
	}

	tx := as.BeginTransaction(as.UserRepository.GetDB())

	if tx.Error != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when start transaction : "+tx.Error.Error())
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

	var user model.User

	if _, err = as.UserRepository.FindByEmailAddress(userInfo.Email); err != nil {
		if err == gorm.ErrRecordNotFound {
			user = model.User{
				Model:                         gorm.Model{},
				Username:                      userInfo.Name,
				Password:                      "",
				Email:                         userInfo.Email,
				EmailVerificationStatus:       true,
				PhoneNumber:                   "",
				PhoneNumberVerificationStatus: false,
				RegistrationMethodID:          2,
			}

			err := as.UserRepository.Create(&user, tx)
			if err != nil {
				tx.Rollback()
				return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when create user : "+err.Error())
			}
		}

		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when get user info : "+err.Error())
	}

	if err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when get user info : "+err.Error())
	}

	accessToken, err := as.createAccessTokens(user.ID)

	if err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when save access token : "+err.Error())
	}

	refreshToken, err := as.createRefreshToken(user.ID)

	if err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when save refresh token : "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when commit transaction : "+err.Error())
	}

	return dto.JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as AuthServiceImplement) ManualRegister(request dto.ManualRegisterRequest) error {
	tx := as.UserRepository.GetDB().Begin()

	if tx.Error != nil {
		tx.Rollback()
		return dto.NewApiError(http.StatusInternalServerError, "Error when start transaction : "+tx.Error.Error())
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

	emailFind, err := as.UserRepository.FindByEmailAddress(request.Email)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return dto.NewApiError(http.StatusInternalServerError, "Error when find email : "+err.Error())
		}
	}

	if emailFind.ID != 0 {
		tx.Rollback()
		return dto.NewApiError(http.StatusBadRequest, "Email already exists")
	}

	if !emailFind.EmailVerificationStatus && emailFind.ID != 0 {
		tx.Rollback()
		return dto.NewApiError(http.StatusBadRequest, "Email already exists but not verified")
	}

	formattedPhoneNumber, err := utils.FormatPhoneNumber(request.PhoneNumber)

	if err != nil {
		tx.Rollback()
		return dto.NewApiError(http.StatusBadRequest, "Error when format phone number : "+err.Error())
	}

	phoneFind, err := as.UserRepository.FindByPhoneNumber(formattedPhoneNumber)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return dto.NewApiError(http.StatusInternalServerError, "Error when find phone number : "+err.Error())
		}
	}

	if phoneFind.ID != 0 {
		tx.Rollback()
		return dto.NewApiError(http.StatusBadRequest, "Phone number already exists")
	}

	if !phoneFind.PhoneNumberVerificationStatus && phoneFind.ID != 0 {
		tx.Rollback()
		return dto.NewApiError(http.StatusBadRequest, "Phone number already exists but not verified")
	}

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		tx.Rollback()
		return dto.NewApiError(http.StatusInternalServerError, "Error when hash password : "+err.Error())
	}

	user := model.User{
		Model:                         gorm.Model{},
		Username:                      request.Name,
		Password:                      hashedPassword,
		Email:                         request.Email,
		EmailVerificationStatus:       false,
		PhoneNumber:                   formattedPhoneNumber,
		PhoneNumberVerificationStatus: false,
		RegistrationMethodID:          1,
	}

	err = as.UserRepository.Create(&user, tx)
	if err != nil {
		tx.Rollback()
		return dto.NewApiError(http.StatusInternalServerError, "Error when create user : "+err.Error())
	}

	tx.SavePoint("sp_after_create_user")

	verifToken, err := utils.GenerateActivationToken()

	if err != nil {
		tx.RollbackTo("sp_after_create_user")
		return dto.NewApiError(http.StatusInternalServerError, "Error when generate activation token : "+err.Error())
	}

	err = as.TokenRepository.CreateVerificationToken(model.VerificationToken{
		UserID: user.ID,
		Token:  verifToken,
	}, tx)

	if err != nil {
		tx.RollbackTo("sp_after_create_user")
		return dto.NewApiError(http.StatusInternalServerError, "Error when save token : "+err.Error())
	}

	err = as.WaService.WASendVerificationLink(formattedPhoneNumber, user.Username, verifToken)

	if err != nil {
		tx.RollbackTo("sp_after_create_user")
		return dto.NewApiError(http.StatusInternalServerError, "Error when send wa verification link : "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.RollbackTo("sp_after_create_user")
		return dto.NewApiError(http.StatusInternalServerError, "Error when commit transaction : "+err.Error())
	}

	return nil
}

func (as AuthServiceImplement) ManualLogin(request dto.ManualLoginRequest) (dto.JwtToken, error) {
	var user model.User

	formattedPhone, err := utils.FormatPhoneNumber(request.PhoneNumber)

	if err != nil {
		return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Error when format phone number : "+err.Error())
	}

	user, err = as.UserRepository.FindByPhoneNumber(formattedPhone)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Phone number or password is wrong")
		}
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when find phone number : "+err.Error())
	}

	if !user.PhoneNumberVerificationStatus {
		return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Phone number not verified")
	}

	if !utils.CheckPassword(user.Password, request.Password) {
		return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Phone number or password is wrong")
	}

	tx := as.BeginTransaction(as.UserRepository.GetDB())
	if tx.Error != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when start transaction : "+tx.Error.Error())
	}

	accessToken, err := as.createAccessTokens(user.ID)

	if err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when save access token : "+err.Error())
	}

	refreshToken, err := as.createRefreshToken(user.ID)
	if err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when save refresh token : "+err.Error())
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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when commit transaction : "+err.Error())
	}

	return dto.JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (as AuthServiceImplement) SendVerificationLink(request dto.SendVerificationLinkRequest) error {
	if request.Type == dto.PhoneVerification {
		formattedPhoneNumber, err := utils.FormatPhoneNumber(request.Value)
		if err != nil {
			return dto.NewApiError(http.StatusBadRequest, "Error when format phone number : "+err.Error())
		}

		user, err := as.UserRepository.FindByPhoneNumber(formattedPhoneNumber)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.NewApiError(http.StatusBadRequest, "Phone number not found")
			}
			return dto.NewApiError(http.StatusInternalServerError, "Error when find phone number : "+err.Error())
		}

		if user.PhoneNumberVerificationStatus {
			return dto.NewApiError(http.StatusBadRequest, "Phone number already verified")
		}

		lastVerifToken, err := as.TokenRepository.FindLastVerificationToken(user.ID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.NewApiError(http.StatusInternalServerError, "Error when find latest token : "+err.Error())
			}
		}

		isTokenExp := lastVerifToken.IsExpired()

		if !isTokenExp {
			return dto.NewApiError(http.StatusBadRequest, "You have to wait 60 minutes to request new token")
		}

		newVerifToken, err := utils.GenerateActivationToken()
		if err != nil {
			return dto.NewApiError(http.StatusInternalServerError, "Error when generate activation token : "+err.Error())
		}

		tx := as.BeginTransaction(as.UserRepository.GetDB())

		if tx.Error != nil {
			tx.Rollback()
			return dto.NewApiError(http.StatusInternalServerError, "Error when start transaction : "+tx.Error.Error())
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

		err = as.TokenRepository.CreateVerificationToken(model.VerificationToken{
			UserID: user.ID,
			Token:  newVerifToken,
		}, tx)

		if err != nil {
			tx.Rollback()
			return dto.NewApiError(http.StatusInternalServerError, "Error when save token : "+err.Error())
		}

		err = as.WaService.WASendVerificationLink(formattedPhoneNumber, user.Username, newVerifToken)
		if err != nil {
			tx.Rollback()
			return dto.NewApiError(http.StatusInternalServerError, "Error when send wa verification link : "+err.Error())
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return dto.NewApiError(http.StatusInternalServerError, "Error when commit transaction : "+err.Error())
		}

		return nil
	}

	return dto.NewApiError(http.StatusBadRequest, "Invalid type")
}

func (as AuthServiceImplement) RefreshToken(request dto.RefreshTokenRequest) (dto.JwtToken, error) {
	_, err := as.TokenRepository.FindRefreshToken(request.RefreshToken)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenLog, err := as.TokenRepository.FindRefreshTokenLog(request.RefreshToken)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when find refresh token log : "+err.Error())
			}
			if tokenLog.ID == 0 {
				return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Token not found")
			}

			return dto.JwtToken{}, dto.NewApiError(http.StatusUnauthorized, "Your account has been used on another device")
		}
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when find token : "+err.Error())
	}

	verif, err := utils.VerifyToken(request.RefreshToken, true)

	if err != nil || !verif.Valid {
		return dto.JwtToken{}, dto.NewApiError(http.StatusBadRequest, "Invalid token")
	}

	claims := verif.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	accessToken, err := utils.GenerateAccessToken(userID)

	if err != nil {
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when generate access token : "+err.Error())
	}

	tx := as.BeginTransaction(as.UserRepository.GetDB())
	if tx.Error != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when start transaction : "+tx.Error.Error())
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

	accessToken, err = as.createAccessTokens(userID)
	if err != nil {
		tx.Rollback()
		return dto.JwtToken{}, dto.NewApiError(http.StatusInternalServerError, "Error when save access token : "+err.Error())
	}

	return dto.JwtToken{
		AccessToken:  accessToken,
		RefreshToken: request.RefreshToken,
	}, nil
}

func (as AuthServiceImplement) createAccessTokens(userID uint) (string, error) {
	accessToken, err := utils.GenerateAccessToken(userID)

	if err != nil {
		return "", errors.New("Error when generate access token : " + err.Error())
	}

	findToken, err := as.TokenRepository.FindAccessTokenByUserId(userID)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("Error when find access token : " + err.Error())
		}
	}

	if findToken.ID != 0 {
		err = as.TokenRepository.UpdateAccessToken(model.AccessToken{
			ID:     findToken.ID,
			Token:  accessToken,
			UserID: userID,
		}, as.tx)
		if err != nil {
			return "", errors.New("Error when update access token : " + err.Error())
		}
	} else {
		err = as.TokenRepository.CreateAccessToken(model.AccessToken{
			Token:  accessToken,
			UserID: userID,
		}, as.tx)

		if err != nil {
			return "", errors.New("Error when save access token : " + err.Error())
		}
	}

	err = as.TokenRepository.CreateAccessTokenLog(model.AccessTokenLogs{
		UserID: userID,
		Token:  accessToken,
	}, as.tx)

	if err != nil {
		return "", errors.New("Error when save access token log : " + err.Error())
	}

	return accessToken, nil
}

func (as AuthServiceImplement) createRefreshToken(userID uint) (string, error) {
	refreshToken, err := utils.GenerateRefreshToken(userID)

	if err != nil {
		return "", errors.New("Error when generate refresh token : " + err.Error())
	}

	findToken, err := as.TokenRepository.FindRefreshTokenByUserId(userID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("Error when find refresh token : " + err.Error())
	}

	if findToken.ID != 0 {
		err = as.TokenRepository.UpdateRefreshToken(model.RefreshToken{
			ID:     findToken.ID,
			Token:  refreshToken,
			UserID: userID},
			as.tx)

		if err != nil {
			return "", errors.New("Error when update refresh token : " + err.Error())
		}
	} else {
		err = as.TokenRepository.CreateRefreshToken(model.RefreshToken{
			Token:  refreshToken,
			UserID: userID,
		}, as.tx)

		if err != nil {
			return "", errors.New("Error when save refresh token : " + err.Error())
		}
	}

	err = as.TokenRepository.CreateRefreshTokenLog(model.RefreshTokenLogs{
		UserID: userID,
		Token:  refreshToken,
	}, as.tx)

	if err != nil {
		return "", errors.New("Error when save refresh token log : " + err.Error())
	}

	return refreshToken, nil
}
