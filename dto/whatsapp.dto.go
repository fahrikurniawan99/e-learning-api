package dto

type ExternalSendVerificationLinkResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
