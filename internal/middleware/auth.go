package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"sso-web-app/internal/models"
	"sso-web-app/internal/services"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return gin.HandlerFunc(func(c *gin.Context) {
		// Try to get token from header
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Try to get token from cookie
			cookie, err := c.Cookie("jwt")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
				c.Abort()
				return
			}
			tokenString = cookie
		}

		// Validate token
		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Get user from database
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)

		c.Next()
	})
}

// OptionalAuthMiddleware checks for authentication but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return gin.HandlerFunc(func(c *gin.Context) {
		// Try to get token from header
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Try to get token from cookie
			cookie, err := c.Cookie("jwt")
			if err != nil {
				c.Next()
				return
			}
			tokenString = cookie
		}

		// Validate token
		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Get user from database
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			c.Next()
			return
		}

		if !user.IsActive {
			c.Next()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)

		c.Next()
	})
}

// GetUserFromContext extracts user from Gin context
func GetUserFromContext(c *gin.Context) *models.User {
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*models.User); ok {
			return u
		}
	}
	return nil
}

// RequireVerified middleware ensures user is verified
func RequireVerified() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		user := GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Email verification required"})
			c.Abort()
			return
		}

		c.Next()
	})
}

// CORS middleware for handling cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
