package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"otp_service/internal/config"
	"otp_service/internal/handler"
	"otp_service/internal/mailer"
	"otp_service/internal/repository"
	"otp_service/internal/service"
)

func main() {
	// 1. Load configuration from .env / environment variables
	cfg := config.LoadConfig()

	// 2. Initialize layers (Dependency Injection)
	repo := repository.NewOtpRepository()
	mail := mailer.NewMailer(cfg)
	otpService := service.NewOtpService(cfg, repo, mail)
	otpHandler := handler.NewOtpHandler(otpService)

	// 3. Setup Fiber app & middleware
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// 4. Register Routes
	// Welcome / API Documentation endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "online",
			"description": "Multi-App Centralized OTP Service",
			"endpoints": fiber.Map{
				"health": "GET /health",
				"send_otp": fiber.Map{
					"methods": []string{"POST /api/send-otp", "POST /send-otp"},
					"body": fiber.Map{
						"email":    "user@example.com (required)",
						"app_name": "My App Name (optional, defaults to 'Your App')",
					},
				},
				"verify_otp": fiber.Map{
					"methods": []string{"POST /api/verify-otp", "POST /verify-otp"},
					"body": fiber.Map{
						"email":    "user@example.com (required)",
						"code":     "123456 (required)",
						"app_name": "My App Name (optional, match the send request)",
					},
				},
			},
		})
	})

	// Health check endpoints
	healthCheck := func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	}
	app.Get("/health", healthCheck)
	app.Get("/helath", healthCheck) // Typo alias

	// Support both /api/* and root paths /*
	api := app.Group("/api")
	api.Post("/send-otp", otpHandler.SendOTP)
	api.Post("/verify-otp", otpHandler.VerifyOTP)

	app.Post("/send-otp", otpHandler.SendOTP)
	app.Post("/verify-otp", otpHandler.VerifyOTP)

	// 5. Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Centralized Multi-App OTP Service starting on port %s...", cfg.Port)
	log.Fatal(app.Listen(addr))
}
