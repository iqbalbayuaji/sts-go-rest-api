package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"recipe-api/models"
	"recipe-api/storage"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// AuthService handles authentication operations
type AuthService struct {
	config       *models.Config
	userStorage  storage.UserStorage
	activeTokens map[string]*TokenInfo
	mutex        sync.RWMutex
}

// TokenInfo stores information about an active token
type TokenInfo struct {
	Username  string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewAuthService creates a new authentication service
func NewAuthService(configPath string, userStorage storage.UserStorage) (*AuthService, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return &AuthService{
		config:       config,
		userStorage:  userStorage,
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
func (as *AuthService) ValidateCredentials(username, password string) (*models.User, bool) {
	user, err := as.userStorage.ValidateCredentials(username, password)
	if err != nil {
		return nil, false
	}
	return user, true
}

// GenerateToken creates a new authentication token
func (as *AuthService) GenerateToken(user *models.User) (string, error) {
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
		Username:  user.Username,
		UserID:    user.ID,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}
	
	return token, nil
}

// ValidateToken checks if a token is valid and not expired
func (as *AuthService) ValidateToken(token string) (*TokenInfo, bool) {
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	
	tokenInfo, exists := as.activeTokens[token]
	if !exists {
		return nil, false
	}
	
	// Check if token is expired
	if time.Now().After(tokenInfo.ExpiresAt) {
		// Remove expired token
		go as.removeToken(token)
		return nil, false
	}
	
	return tokenInfo, true
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
