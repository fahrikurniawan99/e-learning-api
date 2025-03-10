package service

import (
	"elearning_api/config"
	"elearning_api/dto"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GoogleService interface {
	GoogleVerifyToken(code string) (dto.GoogleToken, error)
	GetGoogleUserInfo(idToken string) (*dto.GoogleUser, error)
}

type GoogleServiceImplement struct {
}

func NewGoogleService() GoogleService {
	return &GoogleServiceImplement{}
}

func (gs *GoogleServiceImplement) GoogleVerifyToken(code string) (dto.GoogleToken, error) {
	// Buat form-urlencoded payload
	data := url.Values{}
	data.Set("client_id", config.Env.GoogleClientID)
	data.Set("client_secret", config.Env.GoogleClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", config.Env.GoogleRedirectUri)

	// Kirim request ke Google OAuth API
	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return dto.GoogleToken{}, err
	}
	defer resp.Body.Close()

	// Cek status response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return dto.GoogleToken{}, errors.New("error verifying oauth code: " + string(body))
	}

	// Decode response ke struct GoogleToken
	var token dto.GoogleToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return dto.GoogleToken{}, err
	}

	return token, nil
}

func (gs *GoogleServiceImplement) GetGoogleUserInfo(idToken string) (*dto.GoogleUser, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	var user dto.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
