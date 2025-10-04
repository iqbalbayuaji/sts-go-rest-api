package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recipe-api/auth"
	"recipe-api/models"
	"strings"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// HandleLogin processes login requests
func (ah *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		ah.sendLoginError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		ah.sendLoginError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if loginReq.Username == "" || loginReq.Password == "" {
		ah.sendLoginError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Validate credentials
	user, valid := ah.authService.ValidateCredentials(loginReq.Username, loginReq.Password)
	if !valid {
		ah.sendLoginError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := ah.authService.GenerateToken(user)
	if err != nil {
		ah.sendLoginError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send success response
	response := models.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleLogout processes logout requests
func (ah *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		ah.sendLoginError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header
	token := ah.extractTokenFromHeader(r)
	if token == "" {
		ah.sendLoginError(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	// Validate token exists
	_, valid := ah.authService.ValidateToken(token)
	if !valid {
		ah.sendLoginError(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Invalidate token
	if ah.authService.InvalidateToken(token) {
		response := models.LoginResponse{
			Success: true,
			Message: "Logout successful",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		ah.sendLoginError(w, "Failed to logout", http.StatusInternalServerError)
	}
}

// AuthMiddleware validates authentication for protected routes
func (ah *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Extract token from Authorization header
		token := ah.extractTokenFromHeader(r)
		if token == "" {
			ah.sendLoginError(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Validate token
		tokenInfo, valid := ah.authService.ValidateToken(token)
		if !valid {
			ah.sendLoginError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context (optional, for logging)
		r.Header.Set("X-Username", tokenInfo.Username)
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", tokenInfo.UserID))

		// Call next handler
		next(w, r)
	}
}

// extractTokenFromHeader extracts Bearer token from Authorization header
func (ah *AuthHandler) extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for Bearer token format
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

// sendLoginError sends an authentication error response
func (ah *AuthHandler) sendLoginError(w http.ResponseWriter, message string, statusCode int) {
	response := models.LoginResponse{
		Success: false,
		Error:   message,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
