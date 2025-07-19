package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"sso-web-app/internal/models"
	"sso-web-app/internal/repository"
)

func main() {
	// Initialize repository
	userRepo := repository.NewUserRepository()
	
	// Check if admin user already exists
	existingUsers, err := userRepo.List(1, 0)
	if err == nil && len(existingUsers) > 0 {
		for _, user := range existingUsers {
			if user.IsAdmin || user.Email == "admin@example.com" {
				log.Println("Admin user already exists:", user.Email)
				return
			}
		}
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	
	// Create admin user
	admin := &models.User{
		FirstName:    "System",
		LastName:     "Administrator",
		Email:        "admin@example.com",
		Password:     string(hashedPassword),
		IsActive:     true,
		IsVerified:   true,
		IsAdmin:      true,
		Role:         "admin",
		Bio:          "System administrator account for managing the SSO application",
		Location:     "System",
	}
	
	// Save to database
	createdAdmin, err := userRepo.Create(admin)
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}
	
	log.Printf("Admin user created successfully!")
	log.Printf("Email: %s", createdAdmin.Email)
	log.Printf("Password: admin123")
	log.Printf("User ID: %d", createdAdmin.ID)
	
	// Create some test users
	testUsers := []*models.User{
		{
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
			Role:       "user",
			Bio:        "Regular user account for testing",
			Location:   "New York, USA",
		},
		{
			FirstName:  "Jane",
			LastName:   "Smith",
			Email:      "jane.smith@example.com",
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: false,
			Role:       "moderator",
			Bio:        "Moderator account for testing",
			Location:   "Los Angeles, USA",
		},
		{
			FirstName:  "Bob",
			LastName:   "Johnson",
			Email:      "bob.johnson@example.com",
			Password:   string(hashedPassword),
			IsActive:   false,
			IsVerified: true,
			Role:       "user",
			Bio:        "Inactive user account for testing",
			Location:   "Chicago, USA",
		},
	}
	
	for _, user := range testUsers {
		// Check if user already exists
		existingUser, err := userRepo.GetByEmail(user.Email)
		if err != nil || existingUser == nil {
			// User doesn't exist, create it
			createdUser, err := userRepo.Create(user)
			if err != nil {
				log.Printf("Failed to create test user %s: %v", user.Email, err)
			} else {
				log.Printf("Test user created: %s (ID: %d)", createdUser.Email, createdUser.ID)
			}
		} else {
			log.Printf("Test user already exists: %s", user.Email)
		}
	}
	
	log.Println("Database seeding completed!")
}
