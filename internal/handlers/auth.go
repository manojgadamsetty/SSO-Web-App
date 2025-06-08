package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"sso-web-app/internal/middleware"
	"sso-web-app/internal/models"
	"sso-web-app/internal/services"
)

type AuthHandler struct {
	authService  *services.AuthService
	oauthService *services.OAuthService
}

func NewAuthHandler(authService *services.AuthService, oauthService *services.OAuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
	}
}

// Home renders the home page
func (h *AuthHandler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "SSO Web Application",
	})
}

// LoginPage renders the login page
func (h *AuthHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// RegisterPage renders the registration page
func (h *AuthHandler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set JWT token as HTTP-only cookie
	c.SetCookie("jwt", token, int(time.Hour*24*7/time.Second), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if err == services.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token for the new user
	token, err := h.authService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set JWT token as HTTP-only cookie
	c.SetCookie("jwt", token, int(time.Hour*24*7/time.Second), "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear JWT cookie
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// Dashboard renders the user dashboard
func (h *AuthHandler) Dashboard(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Dashboard",
		"user":  user.ToResponse(),
	})
}

// Profile renders the user profile page
func (h *AuthHandler) Profile(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title": "Profile",
		"user":  user.ToResponse(),
	})
}

// UpdateProfile handles profile updates
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.authService.UpdateProfile(user.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser.ToResponse(),
	})
}

// GetUser returns current user information (API endpoint)
func (h *AuthHandler) GetUser(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}

// UpdateUser handles user updates via API
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.authService.UpdateProfile(user.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser.ToResponse(),
	})
}

// GoogleLogin initiates Google OAuth login
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := h.generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true) // 10 minutes
	
	authURL := h.oauthService.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Verify state parameter
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Handle authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	token, _, err := h.oauthService.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set JWT token as HTTP-only cookie
	c.SetCookie("jwt", token, int(time.Hour*24*7/time.Second), "/", "", false, true)

	// Redirect to dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}

// GitHubLogin initiates GitHub OAuth login
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state := h.generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true) // 10 minutes
	
	authURL := h.oauthService.GetGitHubAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback handles GitHub OAuth callback
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	// Verify state parameter
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Handle authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	token, _, err := h.oauthService.HandleGitHubCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set JWT token as HTTP-only cookie
	c.SetCookie("jwt", token, int(time.Hour*24*7/time.Second), "/", "", false, true)

	// Redirect to dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}

// generateState generates a random state string for OAuth
func (h *AuthHandler) generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
