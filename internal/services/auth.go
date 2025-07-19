package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"sso-web-app/internal/models"
	"sso-web-app/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

// Helper function to convert string to string pointer
func stringPtrAuth(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService() *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}
	
	return &AuthService{
		userRepo:  repository.NewUserRepository(),
		jwtSecret: []byte(secret),
	}
}

// Register creates a new user account
func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	return s.userRepo.Create(user)
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req models.LoginRequest) (string, *models.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// GenerateJWT creates a JWT token for the user
func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT validates a JWT token and returns the user ID
func (s *AuthService) ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, ErrInvalidToken
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		return &models.JWTClaims{
			UserID: uint(userID),
			Email:  email,
		}, nil
	}

	return nil, ErrInvalidToken
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Bio = stringPtrAuth(req.Bio)
	user.Website = stringPtrAuth(req.Website)
	user.Location = stringPtrAuth(req.Location)

	return s.userRepo.Update(user)
}

// HashPassword hashes a plain text password
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// VerifyPassword checks if a password matches the hashed password
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
