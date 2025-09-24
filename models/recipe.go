package models

import (
	"errors"
	"time"
)

// Recipe represents a recipe with all its details
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Ingredients  []string  `json:"ingredients"`
	Instructions string    `json:"instructions"`
	CookingTime  string    `json:"cooking_time"`
	Servings     int       `json:"servings"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate checks if the recipe has all required fields
func (r *Recipe) Validate() error {
	if r.Name == "" {
		return errors.New("recipe name is required")
	}
	if len(r.Ingredients) == 0 {
		return errors.New("at least one ingredient is required")
	}
	if r.Instructions == "" {
		return errors.New("instructions are required")
	}
	if r.CookingTime == "" {
		return errors.New("cooking time is required")
	}
	if r.Servings <= 0 {
		return errors.New("servings must be greater than 0")
	}
	if r.Category == "" {
		return errors.New("category is required")
	}
	return nil
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
