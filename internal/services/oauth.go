package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"sso-web-app/internal/models"
	"sso-web-app/internal/repository"
)

// Helper function to convert string to string pointer
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

type OAuthService struct {
	userRepo     repository.UserRepository
	authService  *AuthService
	googleConfig *oauth2.Config
	githubConfig *oauth2.Config
}

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Given   string `json:"given_name"`
	Family  string `json:"family_name"`
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Blog      string `json:"blog"`
}

func NewOAuthService() *OAuthService {
	googleConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	githubConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	return &OAuthService{
		userRepo:     repository.NewUserRepository(),
		authService:  NewAuthService(),
		googleConfig: googleConfig,
		githubConfig: githubConfig,
	}
}

// GetGoogleAuthURL generates the Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetGitHubAuthURL generates the GitHub OAuth authorization URL
func (s *OAuthService) GetGitHubAuthURL(state string) string {
	return s.githubConfig.AuthCodeURL(state)
}

// HandleGoogleCallback handles the Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(code string) (string, *models.User, error) {
	// Exchange code for token
	token, err := s.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	// Get user info
	googleUser, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Find or create user
	user, err := s.findOrCreateGoogleUser(googleUser)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find or create user: %v", err)
	}

	// Generate JWT token
	jwtToken, err := s.authService.GenerateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT: %v", err)
	}

	return jwtToken, user, nil
}

// HandleGitHubCallback handles the GitHub OAuth callback
func (s *OAuthService) HandleGitHubCallback(code string) (string, *models.User, error) {
	// Exchange code for token
	token, err := s.githubConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	// Get user info
	githubUser, err := s.getGitHubUserInfo(token.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Find or create user
	user, err := s.findOrCreateGitHubUser(githubUser)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find or create user: %v", err)
	}

	// Generate JWT token
	jwtToken, err := s.authService.GenerateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT: %v", err)
	}

	return jwtToken, user, nil
}

func (s *OAuthService) getGoogleUserInfo(accessToken string) (*GoogleUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func (s *OAuthService) getGitHubUserInfo(accessToken string) (*GitHubUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// Get user's primary email if not public
	if githubUser.Email == "" {
		email, err := s.getGitHubUserEmail(accessToken)
		if err == nil {
			githubUser.Email = email
		}
	}

	return &githubUser, nil
}

func (s *OAuthService) getGitHubUserEmail(accessToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary email found")
}

func (s *OAuthService) findOrCreateGoogleUser(googleUser *GoogleUser) (*models.User, error) {
	// Try to find user by Google ID
	user, err := s.userRepo.GetByGoogleID(googleUser.ID)
	if err == nil {
		return user, nil
	}

	// Try to find user by email
	user, err = s.userRepo.GetByEmail(googleUser.Email)
	if err == nil {
		// Update Google ID for existing user
		user.GoogleID = stringPtr(googleUser.ID)
		if user.AvatarURL == nil || *user.AvatarURL == "" {
			user.AvatarURL = stringPtr(googleUser.Picture)
		}
		return s.userRepo.Update(user)
	}

	// Create new user
	user = &models.User{
		Email:     googleUser.Email,
		FirstName: googleUser.Given,
		LastName:  googleUser.Family,
		GoogleID:  stringPtr(googleUser.ID),
		AvatarURL: stringPtr(googleUser.Picture),
		IsActive:  true,
		IsVerified: true, // OAuth users are considered verified
	}

	return s.userRepo.Create(user)
}

func (s *OAuthService) findOrCreateGitHubUser(githubUser *GitHubUser) (*models.User, error) {
	githubIDStr := fmt.Sprintf("%d", githubUser.ID)
	
	// Try to find user by GitHub ID
	user, err := s.userRepo.GetByGitHubID(githubIDStr)
	if err == nil {
		return user, nil
	}

	// Try to find user by email if available
	if githubUser.Email != "" {
		user, err = s.userRepo.GetByEmail(githubUser.Email)
		if err == nil {
			// Update GitHub ID for existing user
			user.GitHubID = stringPtr(githubIDStr)
			if user.AvatarURL == nil || *user.AvatarURL == "" {
				user.AvatarURL = stringPtr(githubUser.AvatarURL)
			}
			return s.userRepo.Update(user)
		}
	}

	// Parse name
	firstName := githubUser.Login
	lastName := ""
	if githubUser.Name != "" {
		firstName = githubUser.Name
	}

	// Create new user
	user = &models.User{
		Email:     githubUser.Email,
		FirstName: firstName,
		LastName:  lastName,
		GitHubID:  stringPtr(githubIDStr),
		AvatarURL: stringPtr(githubUser.AvatarURL),
		Bio:       stringPtr(githubUser.Bio),
		Website:   stringPtr(githubUser.Blog),
		Location:  stringPtr(githubUser.Location),
		IsActive:  true,
		IsVerified: githubUser.Email != "", // Only verified if we have an email
	}

	return s.userRepo.Create(user)
}
