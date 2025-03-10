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

type WhatsappService interface {
	WASendVerificationLink(phoneNumber string, name string, token string) error
}

type WhatsappServiceImplement struct {
}

func NewWhatsappService() WhatsappService {
	return &WhatsappServiceImplement{}
}

func (ws WhatsappServiceImplement) WASendVerificationLink(phoneNumber string, name string, token string) error {
	data := url.Values{}
	data.Set("number", phoneNumber)
	data.Set("name", name)
	data.Set("token", token)

	resp, err := http.Post(config.Env.WaWebUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return errors.New("error when send wa verification link, please try again")
	}
	defer resp.Body.Close()

	// Cek status response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New("error when send wa verification link, please try again : " + string(body))
	}

	// Decode response ke struct GoogleToken
	var responseStruct dto.ExternalSendVerificationLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseStruct); err != nil {
		return errors.New("error when decode response : " + err.Error())
	}

	return nil
}
