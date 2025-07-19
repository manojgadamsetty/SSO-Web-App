package repository

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sso-web-app/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetByGitHubID(githubID string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uint) error
	List(limit, offset int) ([]*models.User, error)
	GetUserStats() (*models.UserStatsResponse, error)
	GetUsersByRole(role string, limit, offset int) ([]*models.User, error)
	SearchUsers(query string, limit, offset int) ([]*models.User, error)
	GetRecentUsers(days int, limit, offset int) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

var db *gorm.DB

func init() {
	var err error
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "sso_app.db"
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Auto migrate the schema
	db.AutoMigrate(&models.User{})
}

func NewUserRepository() UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGitHubID(githubID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("git_hub_id = ?", githubID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetDB returns the database instance for migrations or direct queries
func GetDB() *gorm.DB {
	return db
}

// GetUserStats returns user statistics for admin dashboard
func (r *userRepository) GetUserStats() (*models.UserStatsResponse, error) {
	var stats models.UserStatsResponse
	
	// Total users
	r.db.Model(&models.User{}).Count(&stats.TotalUsers)
	
	// Active users
	r.db.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	
	// Verified users
	r.db.Model(&models.User{}).Where("is_verified = ?", true).Count(&stats.VerifiedUsers)
	
	// Admin users
	r.db.Model(&models.User{}).Where("is_admin = ?", true).Count(&stats.AdminUsers)
	
	// New users today
	r.db.Model(&models.User{}).Where("DATE(created_at) = DATE('now')").Count(&stats.NewUsersToday)
	
	// New users this week
	r.db.Model(&models.User{}).Where("created_at >= datetime('now', '-7 days')").Count(&stats.NewUsersWeek)
	
	// New users this month
	r.db.Model(&models.User{}).Where("created_at >= datetime('now', '-30 days')").Count(&stats.NewUsersMonth)
	
	return &stats, nil
}

// GetUsersByRole returns users filtered by role
func (r *userRepository) GetUsersByRole(role string, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Where("role = ?", role).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// SearchUsers searches users by name or email
func (r *userRepository) SearchUsers(query string, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	searchPattern := "%" + query + "%"
	if err := r.db.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", 
		searchPattern, searchPattern, searchPattern).
		Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetRecentUsers returns users created within the specified number of days
func (r *userRepository) GetRecentUsers(days int, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Where("created_at >= datetime('now', '-? days')", days).
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
