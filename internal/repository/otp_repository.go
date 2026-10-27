package repository

import (
	"strings"
	"sync"
	"time"
)

// OtpData represents the stored OTP and its expiration timestamp
type OtpData struct {
	Code      string
	ExpiresAt time.Time
}

// OtpRepository manages thread-safe in-memory storage of OTPs
type OtpRepository struct {
	mu    sync.RWMutex
	store map[string]OtpData
}

// NewOtpRepository initializes a new in-memory OTP repository
func NewOtpRepository() *OtpRepository {
	return &OtpRepository{
		store: make(map[string]OtpData),
	}
}

// storeKey creates an isolated key per app and email so multiple apps don't collide
func storeKey(email, appName string) string {
	cleanApp := strings.ToLower(strings.TrimSpace(appName))
	if cleanApp == "" {
		cleanApp = "default"
	}
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	return cleanApp + ":" + cleanEmail
}

// Save stores an OTP with an expiration duration for a specific app and email
func (r *OtpRepository) Save(email, appName, code string, ttl time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := storeKey(email, appName)
	r.store[key] = OtpData{
		Code:      code,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves the stored OTP data for an app and email
func (r *OtpRepository) Get(email, appName string) (OtpData, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := storeKey(email, appName)
	data, exists := r.store[key]
	return data, exists
}

// Delete removes the OTP entry after verification
func (r *OtpRepository) Delete(email, appName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := storeKey(email, appName)
	delete(r.store, key)
}
