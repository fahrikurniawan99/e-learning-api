package utils

import (
	"errors"
	"regexp"
	"strings"
)

// IsValidIndonesianPhoneNumber memvalidasi nomor telepon Indonesia.
func IsValidIndonesianPhoneNumber(number string) bool {
	regex := regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)
	return regex.MatchString(number)
}

// FormatPhoneNumber memformat nomor telepon yang valid menjadi format dengan awalan "62".
func FormatPhoneNumber(phoneNumber string) (string, error) {
	if !IsValidIndonesianPhoneNumber(phoneNumber) {
		return "", errors.New("invalid phone number")
	}

	// Hapus semua karakter non-digit
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	switch {
	case strings.HasPrefix(cleaned, "0"):
		return "62" + cleaned[1:], nil
	case strings.HasPrefix(cleaned, "62"):
		return cleaned, nil
	case strings.HasPrefix(cleaned, "8"):
		return "62" + cleaned, nil
	default:
		return "", errors.New("unrecognized phone number format")
	}
}

// IsValidEmail memvalidasi email.
func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
