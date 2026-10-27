package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"otp_service/internal/service"
)

// OtpHandler handles HTTP requests related to OTP
type OtpHandler struct {
	service *service.OtpService
}

// NewOtpHandler creates a new OtpHandler
func NewOtpHandler(service *service.OtpService) *OtpHandler {
	return &OtpHandler{service: service}
}

// SendOtpRequest is the expected JSON payload for sending an OTP
type SendOtpRequest struct {
	Email   string `json:"email"`
	AppName string `json:"app_name"` // e.g. "Jam Room", "FitTrack", "E-Store"
}

// VerifyOtpRequest is the expected JSON payload for verifying an OTP
type VerifyOtpRequest struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	AppName string `json:"app_name"` // optional: must match the app_name passed when sent
}

// SendOTP handles POST /api/send-otp and POST /send-otp
func (h *OtpHandler) SendOTP(c *fiber.Ctx) error {
	var req SendOtpRequest
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Email) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Valid email address is required",
		})
	}

	email := strings.TrimSpace(req.Email)
	appName := strings.TrimSpace(req.AppName)
	if appName == "" {
		appName = "Your App"
	}

	if err := h.service.SendOTP(email, appName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "OTP sent successfully",
		"app_name": appName,
	})
}

// VerifyOTP handles POST /api/verify-otp and POST /verify-otp
func (h *OtpHandler) VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Code) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Both 'email' and 'code' are required",
		})
	}

	email := strings.TrimSpace(req.Email)
	code := strings.TrimSpace(req.Code)
	appName := strings.TrimSpace(req.AppName)
	if appName == "" {
		appName = "Your App"
	}

	isValid, msg := h.service.VerifyOTP(email, appName, code)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   msg,
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"message":  msg,
		"app_name": appName,
	})
}
