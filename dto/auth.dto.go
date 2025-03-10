package dto

type Provider string

const (
	GoogleProvider Provider = "google"
	GithubProvider Provider = "github"
)

type SocialLoginRequest struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

type SocialLoginResponse struct {
	Data  JwtToken     `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}

type ManualRegisterRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type ManualRegisterResponse struct {
	Data  Message      `json:"data"`
	Error string       `json:"error" example:""`
	Meta  MetaResponse `json:"meta"`
}

type ManualLoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type ManualLoginResponse struct {
	Data  JwtToken     `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Data  Message      `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}

type GoogleToken struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type GoogleUser struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type SendVericationLinkType string

const (
	EmailVerification SendVericationLinkType = "email"
	PhoneVerification SendVericationLinkType = "phone"
)

type SendVerificationLinkRequest struct {
	Type  SendVericationLinkType `json:"type" binding:"required"`
	Value string                 `json:"value" binding:"required"`
}

type SendVerificationLinkResponse struct {
	Data  Message      `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	Data  JwtToken     `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}
