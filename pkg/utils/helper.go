package utils

import (
	"fmt"
	"strings"
)

// NormalizePhoneNumber converts Indonesian phone format to E.164
// 08123456789 -> +6281234567890
// 6281234567890 -> +6281234567890
// +6281234567890 -> +6281234567890
func NormalizePhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	if strings.HasPrefix(phone, "0") {
		return "+62" + phone[1:]
	}
	if strings.HasPrefix(phone, "62") && !strings.HasPrefix(phone, "+62") {
		return "+" + phone
	}
	if strings.HasPrefix(phone, "+62") {
		return phone
	}
	return phone
}

// ValidatePassword checks password requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hasLetter := false
	hasNumber := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain both letters and numbers")
	}

	return nil
}
