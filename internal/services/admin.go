package services

import (
	"errors"
	"time"

	"sso-web-app/internal/models"
	"sso-web-app/internal/repository"
)

var (
	ErrNotAuthorized = errors.New("user not authorized for this action")
	ErrInvalidRole   = errors.New("invalid role specified")
)

type AdminService struct {
	userRepo repository.UserRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		userRepo: repository.NewUserRepository(),
	}
}

// IsAdmin checks if user has admin privileges
func (s *AdminService) IsAdmin(user *models.User) bool {
	return user.IsAdmin || user.Role == "admin"
}

// GetUserStats returns dashboard statistics
func (s *AdminService) GetUserStats(adminUser *models.User) (*models.UserStatsResponse, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	return s.userRepo.GetUserStats()
}

// GetAllUsers returns paginated list of all users
func (s *AdminService) GetAllUsers(adminUser *models.User, limit, offset int) ([]*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	return s.userRepo.List(limit, offset)
}

// GetUsersByRole returns users filtered by role
func (s *AdminService) GetUsersByRole(adminUser *models.User, role string, limit, offset int) ([]*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	validRoles := map[string]bool{
		"user":      true,
		"admin":     true,
		"moderator": true,
	}
	
	if !validRoles[role] {
		return nil, ErrInvalidRole
	}
	
	return s.userRepo.GetUsersByRole(role, limit, offset)
}

// SearchUsers searches for users by name or email
func (s *AdminService) SearchUsers(adminUser *models.User, query string, limit, offset int) ([]*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	return s.userRepo.SearchUsers(query, limit, offset)
}

// GetRecentUsers returns recently registered users
func (s *AdminService) GetRecentUsers(adminUser *models.User, days, limit, offset int) ([]*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	return s.userRepo.GetRecentUsers(days, limit, offset)
}

// GetUserByID returns a specific user by ID
func (s *AdminService) GetUserByID(adminUser *models.User, userID uint) (*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	return s.userRepo.GetByID(userID)
}

// UpdateUser updates user information (admin operation)
func (s *AdminService) UpdateUser(adminUser *models.User, userID uint, req models.AdminUpdateUserRequest) (*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	// Get the user to update
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	// Prevent non-super-admin from modifying other admins
	if user.IsAdmin && adminUser.Role != "admin" {
		return nil, ErrNotAuthorized
	}
	
	// Update fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Bio = req.Bio
	user.Website = req.Website
	user.Location = req.Location
	
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	
	if req.IsVerified != nil {
		user.IsVerified = *req.IsVerified
	}
	
	if req.IsAdmin != nil {
		// Only super admins can modify admin status
		if adminUser.Role == "admin" {
			user.IsAdmin = *req.IsAdmin
		}
	}
	
	if req.Role != "" {
		validRoles := map[string]bool{
			"user":      true,
			"admin":     true,
			"moderator": true,
		}
		
		if !validRoles[req.Role] {
			return nil, ErrInvalidRole
		}
		
		// Only super admins can assign admin role
		if req.Role == "admin" && adminUser.Role != "admin" {
			return nil, ErrNotAuthorized
		}
		
		user.Role = req.Role
	}
	
	return s.userRepo.Update(user)
}

// DeactivateUser deactivates a user account
func (s *AdminService) DeactivateUser(adminUser *models.User, userID uint) (*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	// Prevent deactivating other admins unless super admin
	if user.IsAdmin && adminUser.Role != "admin" {
		return nil, ErrNotAuthorized
	}
	
	// Prevent self-deactivation
	if user.ID == adminUser.ID {
		return nil, errors.New("cannot deactivate your own account")
	}
	
	user.IsActive = false
	return s.userRepo.Update(user)
}

// ActivateUser activates a user account
func (s *AdminService) ActivateUser(adminUser *models.User, userID uint) (*models.User, error) {
	if !s.IsAdmin(adminUser) {
		return nil, ErrNotAuthorized
	}
	
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	user.IsActive = true
	return s.userRepo.Update(user)
}

// DeleteUser permanently deletes a user account
func (s *AdminService) DeleteUser(adminUser *models.User, userID uint) error {
	if !s.IsAdmin(adminUser) {
		return ErrNotAuthorized
	}
	
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	
	// Prevent deleting other admins unless super admin
	if user.IsAdmin && adminUser.Role != "admin" {
		return ErrNotAuthorized
	}
	
	// Prevent self-deletion
	if user.ID == adminUser.ID {
		return errors.New("cannot delete your own account")
	}
	
	return s.userRepo.Delete(userID)
}

// PromoteToAdmin promotes a user to admin role
func (s *AdminService) PromoteToAdmin(adminUser *models.User, userID uint) (*models.User, error) {
	if !s.IsAdmin(adminUser) || adminUser.Role != "admin" {
		return nil, ErrNotAuthorized
	}
	
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	user.IsAdmin = true
	user.Role = "admin"
	return s.userRepo.Update(user)
}

// DemoteFromAdmin removes admin privileges from a user
func (s *AdminService) DemoteFromAdmin(adminUser *models.User, userID uint) (*models.User, error) {
	if !s.IsAdmin(adminUser) || adminUser.Role != "admin" {
		return nil, ErrNotAuthorized
	}
	
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	// Prevent self-demotion
	if user.ID == adminUser.ID {
		return nil, errors.New("cannot demote your own account")
	}
	
	user.IsAdmin = false
	user.Role = "user"
	return s.userRepo.Update(user)
}
