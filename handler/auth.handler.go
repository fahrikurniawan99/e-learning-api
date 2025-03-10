package handler

import (
	"elearning_api/dto"
	"elearning_api/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler interface {
	RegisterOrLogin(c *gin.Context)
	ManualRegister(c *gin.Context)
	ManualLogin(c *gin.Context)
	SendVerificationLink(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type AuthHandlerImplement struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &AuthHandlerImplement{
		AuthService: authService,
	}
}

// TO GET OAUTH2 CODE
// https://accounts.google.com/o/oauth2/v2/auth?client_id=393747339679-4udk2p70dvn8sc772cqnr86bne4cechi.apps.googleusercontent.com&redirect_uri=http://localhost:3000/api/auth/callback/google&response_type=code&scope=openid%20email%20profile&access_type=offline&prompt=consent
// https://github.com/login/oauth/authorize?client_id=Ov23lie6Emz6dgOqUlZI&redirect_uri=http://localhost:3000/callback/github&scope=read:user

// @Summary Social Login
// @Description Register atau login dengan Google atau GitHub
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SocialLoginRequest true "Social Login Request"
// @Success 200 {object} dto.SocialLoginResponse
// @Router /auth/social [post]
func (ah AuthHandlerImplement) RegisterOrLogin(c *gin.Context) {
	var request dto.SocialLoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(request)
		c.JSON(http.StatusBadRequest, dto.SocialLoginResponse{
			Error: "Error when bind data : " + err.Error(),
		})
		return
	}

	token, err := ah.AuthService.SocialLogin(request)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.SocialLoginResponse{
				Data:  dto.JwtToken{},
				Error: apiErr.Message,
				Meta:  dto.MetaResponse{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SocialLoginResponse{
			Data:  dto.JwtToken{},
			Error: "Unexpected Error",
			Meta:  dto.MetaResponse{},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SocialLoginResponse{
		Data:  token,
		Error: "",
		Meta:  dto.MetaResponse{},
	})

}

// @Summary Manual Register
// @Description Register dengan email dan password secara manual
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ManualRegisterRequest true "Manual Register Request"
// @Success 200 {object} dto.ManualRegisterResponse
// @Router /auth/register [post]
func (ah AuthHandlerImplement) ManualRegister(c *gin.Context) {
	var request dto.ManualRegisterRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ManualRegisterResponse{
			Error: "Error when bind data : " + err.Error(),
		})
		return
	}

	err = ah.AuthService.ManualRegister(request)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.ManualRegisterResponse{
				Error: apiErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ManualRegisterResponse{
			Error: "Unexpected Error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ManualRegisterResponse{
		Error: "",
		Data:  dto.Message{Message: "Register Success, please your whatsapp to verify your phone number"},
	})
}

// @Summary Manual Login
// @Description Login dengan nomor wa dan password secara manual
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ManualLoginRequest true "Manual Login Request"
// @Success 200 {object} dto.ManualLoginResponse
// @Router /auth/login [post]
func (ah AuthHandlerImplement) ManualLogin(c *gin.Context) {
	request := dto.ManualLoginRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ManualLoginResponse{
			Error: "Error when bind data : " + err.Error(),
		})
		return
	}

	token, err := ah.AuthService.ManualLogin(request)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.ManualLoginResponse{
				Data:  dto.JwtToken{},
				Error: apiErr.Message,
				Meta:  dto.MetaResponse{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ManualLoginResponse{
			Data:  dto.JwtToken{},
			Error: "Unexpected Error",
			Meta:  dto.MetaResponse{},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ManualLoginResponse{
		Data:  token,
		Error: "",
		Meta:  dto.MetaResponse{},
	})
}

// @Summary Send Verification Link
// @Description Kirim link verifikasi ke email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SendVerificationLinkRequest true "User Profile Request"
// @Success 200 {object} dto.SendVerificationLinkResponse
// @Router /auth/verification/send [post]
func (ah AuthHandlerImplement) SendVerificationLink(c *gin.Context) {
	var request dto.SendVerificationLinkRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.SendVerificationLinkResponse{
			Error: "Error when bind data : " + err.Error(),
		})
		return
	}

	err := ah.AuthService.SendVerificationLink(request)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.SendVerificationLinkResponse{
				Error: apiErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SendVerificationLinkResponse{
			Error: "Unexpected Error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SendVerificationLinkResponse{
		Error: "",
		Data:  dto.Message{Message: "Verification Link Sent"},
	})
}

// @Summary Refresh Token
// @Description Refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} dto.RefreshTokenResponse
// @Router /auth/refresh-token [post]
func (ah AuthHandlerImplement) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.RefreshTokenResponse{
			Error: "Error when bind data : " + err.Error(),
		})
		return
	}

	token, err := ah.AuthService.RefreshToken(request)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.RefreshTokenResponse{
				Data:  dto.JwtToken{},
				Error: apiErr.Message,
				Meta:  dto.MetaResponse{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.RefreshTokenResponse{
			Data:  dto.JwtToken{},
			Error: "Unexpected Error",
			Meta:  dto.MetaResponse{},
		})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		Data:  token,
		Error: "",
		Meta:  dto.MetaResponse{},
	})
}
