package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"otp_service/internal/config"
	"otp_service/internal/mailer"
	"otp_service/internal/repository"
)

// OtpService contains business logic for OTP generation, sending, and verification
type OtpService struct {
	cfg    *config.Config
	repo   *repository.OtpRepository
	mailer *mailer.Mailer
}

// NewOtpService creates a new OtpService instance
func NewOtpService(cfg *config.Config, repo *repository.OtpRepository, mailer *mailer.Mailer) *OtpService {
	return &OtpService{
		cfg:    cfg,
		repo:   repo,
		mailer: mailer,
	}
}

// GenerateOTP creates a cryptographically secure 6-digit random code
func (s *OtpService) GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// SendOTP generates an OTP, sends the email customized for appName via Brevo, and saves it in repository
func (s *OtpService) SendOTP(email, appName string) error {
	if strings.TrimSpace(s.cfg.BrevoApiKey) == "" || strings.TrimSpace(s.cfg.BrevoSenderEmail) == "" {
		return errors.New("brevo configuration missing: please set BREVO_API_KEY and BREVO_SENDER_EMAIL in environment")
	}

	cleanApp := strings.TrimSpace(appName)
	if cleanApp == "" {
		cleanApp = "Your App"
	}

	code, err := s.GenerateOTP()
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	// Send email and verify delivery before saving to store
	if err := s.mailer.SendOTP(email, code, cleanApp); err != nil {
		log.Printf("[ERROR] Failed to send OTP email to %s for app '%s': %v\n", email, cleanApp, err)
		return fmt.Errorf("failed to deliver email: %w", err)
	}

	ttl := time.Duration(s.cfg.OTPExpiryMinutes) * time.Minute
	s.repo.Save(email, cleanApp, code, ttl)
	log.Printf("[SUCCESS] OTP email successfully sent to %s for app '%s'\n", email, cleanApp)

	return nil
}

// VerifyOTP validates the code against what's saved in the repository for the given appName
func (s *OtpService) VerifyOTP(email, appName, code string) (bool, string) {
	cleanApp := strings.TrimSpace(appName)
	if cleanApp == "" {
		cleanApp = "Your App"
	}

	data, exists := s.repo.Get(email, cleanApp)
	if !exists {
		return false, fmt.Sprintf("OTP not found for %s. Please request a new one", cleanApp)
	}

	if time.Now().After(data.ExpiresAt) {
		s.repo.Delete(email, cleanApp)
		return false, "OTP has expired. Please request a new one"
	}

	if data.Code != code {
		return false, "Invalid OTP code"
	}

	// Single-use guarantee: remove OTP immediately after successful verification
	s.repo.Delete(email, cleanApp)
	return true, "OTP verified successfully"
}
