package dto

type UserProfileRequest struct {
	UserID uint `json:"user_id"`
}

type UserBasicInfo struct {
	ID                            uint   `json:"id"`
	Username                      string `json:"username"`
	Email                         string `json:"email"`
	PhoneNumber                   string `json:"phone_number"`
	EmailVerificationStatus       bool   `json:"email_verification_status"`
	PhoneNumberVerificationStatus bool   `json:"phone_number_verification_status"`
}

type UserProfileResponse struct {
	Data  UserBasicInfo `json:"data"`
	Error string        `json:"error"`
	Meta  MetaResponse  `json:"meta"`
}

type UserConfirmRequest struct {
	Token string `json:"token"`
}

type UserConfirmResponse struct {
	Data  Message      `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}

type UserUpdateStatusVerifPhone struct {
	ID                            uint `json:"id"`
	PhoneNumberVerificationStatus bool
}
