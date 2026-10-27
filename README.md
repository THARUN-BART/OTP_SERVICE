# Multi-App Centralized OTP Service (Brevo Edition)

A clean, high-performance OTP (One-Time Password) microservice written in **Go (Fiber)** powered exclusively by **Brevo REST API**.

Deploy **once** on Render (or any cloud) and use it across all your applications (Flutter, React, iOS, Android) with:
- **Zero custom domain required**: Works using your verified personal Gmail.
- **300 free emails/day** via Brevo.
- **Port 443 HTTPS only**: Completely avoids cloud firewall SMTP blocks on ports 25, 465, and 587.

---

## 📁 Architecture

```
otp_service/
├── Dockerfile                   # Multi-stage production build (Alpine runtime)
├── .env.example                 # Clean Brevo-only environment template
├── go.mod / go.sum              # Go dependencies
├── main.go                      # Entry point: dependency wiring & Fiber routes
└── internal/
    ├── config/config.go         # Environment loader (Brevo keys & PORT)
    ├── repository/otp_repository.go # Thread-safe in-memory storage (RWMutex)
    ├── mailer/smtp_mailer.go    # Brevo HTTPS REST API client
    ├── service/otp_service.go   # Business logic (OTP generation, TTL & validation)
    └── handler/otp_handler.go   # HTTP endpoints (Fiber route handlers)
```

---

## ⚙️ Environment Variables

All you need are these 3 variables:

```env
# Brevo API key (starts with xkeysib-...)
BREVO_API_KEY=xkeysib-your_brevo_api_key_here

# The verified email address in your Brevo account
BREVO_SENDER_EMAIL=tharunp29112006@gmail.com

# OTP Expiry duration in minutes (default: 5)
OTP_EXPIRY_MINUTES=5
```

---

## 🚀 API Endpoints

### 1. Health Check
```http
GET /health
```
```json
{"status": "healthy"}
```

---

### 2. Send OTP
```http
POST /api/send-otp
Content-Type: application/json
```
```json
{
  "email": "customer@gmail.com",
  "app_name": "Jam Room"
}
```
**Response (`200 OK`)**:
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "app_name": "Jam Room"
}
```

---

### 3. Verify OTP
```http
POST /api/verify-otp
Content-Type: application/json
```
```json
{
  "email": "customer@gmail.com",
  "code": "123456",
  "app_name": "Jam Room"
}
```
**Response (`200 OK`)**:
```json
{
  "success": true,
  "message": "OTP verified successfully",
  "app_name": "Jam Room"
}
```

---

## 🐳 Docker Deployment

```bash
docker build -t bartsimpson2911/otp:latest .
docker push bartsimpson2911/otp:latest
```
