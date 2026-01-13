package whatsapp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
)

// WhatsAppClient handles WhatsApp API communication using Ultramsg
type WhatsAppClient struct {
	InstanceID   string
	APIToken     string
	BaseURL      string
	SenderNumber string
	httpClient   *http.Client
}

// NewWhatsAppClient creates a new WhatsApp client
func NewWhatsAppClient(instanceID, apiToken, baseURL, senderNumber string) *WhatsAppClient {
	return &WhatsAppClient{
		InstanceID:   instanceID,
		APIToken:     apiToken,
		BaseURL:      baseURL,
		SenderNumber: senderNumber,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendOTP sends OTP code via WhatsApp
func (w *WhatsAppClient) SendOTP(phoneNumber, otpCode string) error {
	message := fmt.Sprintf(
		"🔐 *Ojek Kampus - Kode OTP*\n\n"+
			"Kode OTP Anda: *%s*\n\n"+
			"Berlaku selama 5 menit.\n"+
			"Jangan bagikan kode ini kepada siapapun!",
		otpCode,
	)

	return w.SendMessage(phoneNumber, message)
}

// SendMessage sends a WhatsApp message
func (w *WhatsAppClient) SendMessage(phoneNumber, message string) error {
	// Format phone number (remove leading 0, add 62)
	formattedPhone := w.formatPhoneNumber(phoneNumber)

	// Build API URL
	apiURL := fmt.Sprintf("%s/%s/messages/chat", w.BaseURL, w.InstanceID)

	// Prepare form data
	data := url.Values{}
	data.Set("token", w.APIToken)
	data.Set("to", formattedPhone)
	data.Set("body", message)

	// Create request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create WhatsApp request")
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := w.httpClient.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to send WhatsApp message")
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to read WhatsApp response")
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.Log.Error().
			Int("status", resp.StatusCode).
			Str("response", string(body)).
			Msg("WhatsApp API error")
		return fmt.Errorf("WhatsApp API error: status %d, response: %s", resp.StatusCode, string(body))
	}

	logger.Log.Info().
		Str("phone", formattedPhone).
		Str("response", string(body)).
		Msg("WhatsApp message sent successfully")

	return nil
}

// formatPhoneNumber converts Indonesian phone format to international format
// Examples:
//   - 081234567890 -> 6281234567890
//   - +6281234567890 -> 6281234567890
//   - 6281234567890 -> 6281234567890
func (w *WhatsAppClient) formatPhoneNumber(phone string) string {
	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Remove + prefix
	phone = strings.TrimPrefix(phone, "+")

	// If starts with 0, replace with 62
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	// If doesn't start with 62, add it
	if !strings.HasPrefix(phone, "62") {
		phone = "62" + phone
	}

	return phone
}
