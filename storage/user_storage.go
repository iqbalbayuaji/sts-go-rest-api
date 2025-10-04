package storage

import (
	"database/sql"
	"fmt"
	"recipe-api/database"
	"recipe-api/models"

	"golang.org/x/crypto/bcrypt"
)

// PostgresUserStorage handles PostgreSQL operations for users
type PostgresUserStorage struct {
	db *sql.DB
}

// NewPostgresUserStorage creates a new PostgreSQL user storage instance
func NewPostgresUserStorage() *PostgresUserStorage {
	return &PostgresUserStorage{
		db: database.GetDB(),
	}
}

// GetUserByUsername retrieves a user by username
func (pus *PostgresUserStorage) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, email, is_active, created_at, updated_at, created_by, updated_by
		FROM users
		WHERE username = $1 AND is_active = true
	`

	var user models.User
	err := pus.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (pus *PostgresUserStorage) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, email, is_active, created_at, updated_at, created_by, updated_by
		FROM users
		WHERE id = $1 AND is_active = true
	`

	var user models.User
	err := pus.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// CreateUser creates a new user
func (pus *PostgresUserStorage) CreateUser(user models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (username, password_hash, email, is_active, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err = pus.db.QueryRow(
		query,
		user.Username, string(hashedPassword), user.Email, user.IsActive,
		user.CreatedBy, user.UpdatedBy,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// UpdateUser updates an existing user
func (pus *PostgresUserStorage) UpdateUser(user models.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, is_active = $4, updated_by = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := pus.db.QueryRow(
		query,
		user.ID, user.Username, user.Email, user.IsActive, user.UpdatedBy,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// DeleteUser soft deletes a user by setting is_active to false
func (pus *PostgresUserStorage) DeleteUser(id int) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := pus.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// ValidateCredentials validates username and password
func (pus *PostgresUserStorage) ValidateCredentials(username, password string) (*models.User, error) {
	user, err := pus.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// UpdatePassword updates user password
func (pus *PostgresUserStorage) UpdatePassword(userID int, newPassword string, updatedBy *int) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		UPDATE users 
		SET password_hash = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := pus.db.Exec(query, userID, string(hashedPassword), updatedBy)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}
