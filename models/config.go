package models

// Config represents the application configuration
type Config struct {
	Users            []User `yaml:"users"`
	JWTSecret        string `yaml:"jwt_secret"`
	TokenExpiryHours int    `yaml:"token_expiry_hours"`
}

// User represents a user from the config file
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
