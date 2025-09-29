package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"recipe-api/models"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// AuthService handles authentication operations
type AuthService struct {
	config       *models.Config
	activeTokens map[string]*TokenInfo
	mutex        sync.RWMutex
}

// TokenInfo stores information about an active token
type TokenInfo struct {
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewAuthService creates a new authentication service
func NewAuthService(configPath string) (*AuthService, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return &AuthService{
		config:       config,
		activeTokens: make(map[string]*TokenInfo),
	}, nil
}

// loadConfig loads configuration from YAML file
func loadConfig(configPath string) (*models.Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default values if not specified
	if config.TokenExpiryHours == 0 {
		config.TokenExpiryHours = 24
	}
	if config.JWTSecret == "" {
		config.JWTSecret = "default-secret-change-this"
	}

	return &config, nil
}

// ValidateCredentials checks if username and password are valid
func (as *AuthService) ValidateCredentials(username, password string) bool {
	for _, user := range as.config.Users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

// GenerateToken creates a new authentication token
func (as *AuthService) GenerateToken(username string) (string, error) {
	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	
	token := hex.EncodeToString(tokenBytes)
	
	// Store token info
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	now := time.Now()
	expiresAt := now.Add(time.Duration(as.config.TokenExpiryHours) * time.Hour)
	
	as.activeTokens[token] = &TokenInfo{
		Username:  username,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}
	
	return token, nil
}

// ValidateToken checks if a token is valid and not expired
func (as *AuthService) ValidateToken(token string) (string, bool) {
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	
	tokenInfo, exists := as.activeTokens[token]
	if !exists {
		return "", false
	}
	
	// Check if token is expired
	if time.Now().After(tokenInfo.ExpiresAt) {
		// Remove expired token
		go as.removeToken(token)
		return "", false
	}
	
	return tokenInfo.Username, true
}

// InvalidateToken removes a token from active tokens
func (as *AuthService) InvalidateToken(token string) bool {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	_, exists := as.activeTokens[token]
	if exists {
		delete(as.activeTokens, token)
		return true
	}
	return false
}

// removeToken removes a token (used for cleanup)
func (as *AuthService) removeToken(token string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	delete(as.activeTokens, token)
}

// CleanupExpiredTokens removes all expired tokens
func (as *AuthService) CleanupExpiredTokens() {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	now := time.Now()
	for token, tokenInfo := range as.activeTokens {
		if now.After(tokenInfo.ExpiresAt) {
			delete(as.activeTokens, token)
		}
	}
}

// GetActiveTokensCount returns the number of active tokens
func (as *AuthService) GetActiveTokensCount() int {
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	return len(as.activeTokens)
}
