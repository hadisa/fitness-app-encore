// Package admin provides functionality for user management and authentication.
package admin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"encore.dev/storage/sqldb"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserDetail contains additional user information
type UserDetail struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Fullname   string    `json:"fullname"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	ProvinceID *int      `json:"province_id,omitempty"`
	DistrictID *int      `json:"district_id,omitempty"`
	CityID     *int      `json:"city_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AuthResponse is returned after successful authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// LoginParams contains the data needed to authenticate a user
type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ProfileResponse contains the user's profile information
type ProfileResponse struct {
	User       *User       `json:"user"`
	UserDetail *UserDetail `json:"user_detail"`
}

// RegisterParams contains the data needed to register a new user
type RegisterParams struct {
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Fullname   string  `json:"fullname"`
	Address    *string `json:"address,omitempty"`
	ProvinceID *int    `json:"province_id,omitempty"`
	CityID     *int    `json:"city_id,omitempty"`
	DistrictID *int    `json:"district_id,omitempty"`
}

// Register creates a new user account
//
//encore:api public method=POST path=/admin/register
func Register(ctx context.Context, params *RegisterParams) (*AuthResponse, error) {
	// Hash the password
	hashedPassword := hashPassword(params.Password)

	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert user
	var userID int
	err = tx.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`, params.Username, params.Email, hashedPassword).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Insert user details
	_, err = tx.Exec(ctx, `
		INSERT INTO user_details (
			user_id, fullname, address, province_id, city_id, district_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, userID, params.Fullname, params.Address, params.ProvinceID, params.CityID, params.DistrictID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Get the created user
	user := &User{
		ID:        userID,
		Username:  params.Username,
		Email:     params.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate token (in a real app, use JWT or similar)
	token := generateToken(userID)

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user
//
//encore:api public method=POST path=/admin/login
func Login(ctx context.Context, params *LoginParams) (*AuthResponse, error) {
	// Get user by username
	var user User
	var hashedPassword string
	err := db.QueryRow(ctx, `
        SELECT id, username, email, password_hash, created_at, updated_at
        FROM users
        WHERE username = $1
    `, params.Username).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !verifyPassword(params.Password, hashedPassword) {
		return nil, errors.New("invalid credentials")
	}

	// Rest of the function remains the same
	token := generateToken(user.ID)

	return &AuthResponse{
		Token: token,
		User:  &user,
	}, nil
}

// GetProfile returns the current user's profile

func GetProfile(ctx context.Context) (*ProfileResponse, error) {
	// In a real app, get user ID from context (from JWT or session)
	userID := 1 // This should come from the authenticated context

	// Get user
	var user User
	err := db.QueryRow(ctx, `
        SELECT id, username, email, created_at, updated_at
        FROM users
        WHERE id = $1
    `, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get user details
	var detail UserDetail
	err = db.QueryRow(ctx, `
        SELECT id, user_id, fullname, address, postal_code, 
               province_id, city_id, district_id, created_at, updated_at
        FROM user_details
        WHERE user_id = $1
    `, userID).Scan(
		&detail.ID, &detail.UserID, &detail.Fullname, &detail.Address, &detail.PostalCode,
		&detail.ProvinceID, &detail.CityID, &detail.DistrictID, &detail.CreatedAt, &detail.UpdatedAt,
	)
	if err != nil {
		// User details might not exist yet
		detail = UserDetail{UserID: userID}
	}

	return &ProfileResponse{
		User:       &user,
		UserDetail: &detail,
	}, nil
}

// Helper function to hash password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Helper function to verify password
func verifyPassword(password, hashedPassword string) bool {
	return hashPassword(password) == hashedPassword
}

// Helper function to generate token (in a real app, use JWT)
func generateToken(userID int) string {
	// This is a simplified example. In production, use a proper JWT implementation
	return "generated-jwt-token-for-user-" + string(userID)
}

// Define the database connection
var db = sqldb.Named("admin")
