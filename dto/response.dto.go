package dto

type MetaResponse struct {
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
}

type Response struct {
	Data  interface{}  `json:"data"`
	Error string       `json:"error"`
	Meta  MetaResponse `json:"meta"`
}

//type Token struct {
//	Token string `json:"token"`
//}

type JwtToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Code struct {
	Code string `json:"code"`
}

type Message struct {
	Message string `json:"message"`
}

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}
