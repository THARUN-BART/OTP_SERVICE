package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"otp_service/internal/config"
)

// Mailer handles sending emails via Brevo HTTPS REST API
type Mailer struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewMailer creates a new Mailer instance with the given configuration
func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// cleanAppName returns the app name or a default fallback
func cleanAppName(appName string) string {
	trimmed := strings.TrimSpace(appName)
	if trimmed == "" {
		return "Your App"
	}
	return trimmed
}

// SendOTP dispatches the OTP email via Brevo's REST API over HTTPS (Port 443)
func (m *Mailer) SendOTP(to, otp, appName string) error {
	resolvedApp := cleanAppName(appName)
	senderEmail := strings.TrimSpace(m.cfg.BrevoSenderEmail)
	apiKey := strings.TrimSpace(m.cfg.BrevoApiKey)

	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY environment variable is missing")
	}
	if senderEmail == "" {
		return fmt.Errorf("BREVO_SENDER_EMAIL environment variable is missing")
	}

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  resolvedApp,
			"email": senderEmail,
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject": fmt.Sprintf("Your %s Verification Code", resolvedApp),
		"htmlContent": fmt.Sprintf(`
			<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 480px; margin: 0 auto; padding: 24px; border: 1px solid #e2e8f0; border-radius: 12px; background-color: #ffffff;">
				<h2 style="color: #0f172a; margin-top: 0; font-size: 22px;">%s Verification</h2>
				<p style="color: #475569; font-size: 15px; line-height: 1.5;">Hello,</p>
				<p style="color: #475569; font-size: 15px; line-height: 1.5;">Please use the following 6-digit verification code to complete your sign-in:</p>
				<div style="background-color: #f8fafc; border: 1px dashed #cbd5e1; padding: 18px; text-align: center; border-radius: 8px; margin: 24px 0;">
					<span style="font-size: 32px; font-weight: 700; letter-spacing: 8px; color: #1e293b; display: inline-block;">%s</span>
				</div>
				<p style="color: #64748b; font-size: 14px; margin-bottom: 8px;">⏳ This code will expire in <strong>%d minutes</strong>.</p>
				<hr style="border: none; border-top: 1px solid #f1f5f9; margin: 20px 0;" />
				<p style="color: #94a3b8; font-size: 12px; margin: 0;">If you did not request this code, you can safely ignore this email.</p>
			</div>
		`, resolvedApp, otp, m.cfg.OTPExpiryMinutes),
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode Brevo payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create Brevo HTTP request: %w", err)
	}

	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Brevo API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
