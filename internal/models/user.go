package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Password    string `gorm:"not null" json:"-"` // Never include password in JSON
	FirstName   string `gorm:"not null" json:"first_name"`
	LastName    string `gorm:"not null" json:"last_name"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	IsVerified  bool   `gorm:"default:false" json:"is_verified"`
	
	// OAuth fields
	GoogleID  string `gorm:"uniqueIndex" json:"google_id,omitempty"`
	GitHubID  string `gorm:"uniqueIndex" json:"github_id,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	
	// Profile fields
	Bio       string `json:"bio,omitempty"`
	Website   string `json:"website,omitempty"`
	Location  string `json:"location,omitempty"`
	
	// Security fields
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	PasswordResetAt *time.Time `json:"password_reset_at,omitempty"`
}

// UserResponse represents user data returned to clients
type UserResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	Website     string    `json:"website,omitempty"`
	Location    string    `json:"location,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		IsActive:    u.IsActive,
		IsVerified:  u.IsVerified,
		AvatarURL:   u.AvatarURL,
		Bio:         u.Bio,
		Website:     u.Website,
		Location:    u.Location,
		CreatedAt:   u.CreatedAt,
		LastLoginAt: u.LastLoginAt,
	}
}

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
}

// UpdateProfileRequest represents profile update request data
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Bio       string `json:"bio"`
	Website   string `json:"website"`
	Location  string `json:"location"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}
