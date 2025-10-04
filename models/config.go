package models

import "time"

// Config represents the application configuration
type Config struct {
	Database         DatabaseConfig `yaml:"database"`
	JWTSecret        string         `yaml:"jwt_secret"`
	TokenExpiryHours int            `yaml:"token_expiry_hours"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// User represents a user in the database
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password_hash"` // Don't expose password in JSON
	Email     string    `json:"email" db:"email"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *int      `json:"created_by" db:"created_by"`
	UpdatedBy *int      `json:"updated_by" db:"updated_by"`
}
