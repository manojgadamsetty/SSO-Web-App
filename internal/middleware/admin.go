package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sso-web-app/internal/models"
)

// AdminRequired middleware checks if the authenticated user has admin privileges
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		user, exists := c.Get("user")
		if !exists {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"title": "Unauthorized",
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		authUser, ok := user.(*models.User)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"title": "Error",
				"error": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Check if user has admin privileges
		if !authUser.IsAdmin && authUser.Role != "admin" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"title": "Access Denied",
				"error": "Admin privileges required to access this page",
			})
			c.Abort()
			return
		}

		// User is admin, continue
		c.Next()
	}
}

// SuperAdminRequired middleware checks if the authenticated user has super admin privileges
func SuperAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		user, exists := c.Get("user")
		if !exists {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"title": "Unauthorized",
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		authUser, ok := user.(*models.User)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"title": "Error",
				"error": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Check if user has super admin privileges
		if authUser.Role != "admin" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"title": "Access Denied",
				"error": "Super admin privileges required to access this page",
			})
			c.Abort()
			return
		}

		// User is super admin, continue
		c.Next()
	}
}

// AdminAPIRequired middleware for API endpoints that require admin privileges
func AdminAPIRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		authUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Check if user has admin privileges
		if !authUser.IsAdmin && authUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin privileges required",
			})
			c.Abort()
			return
		}

		// User is admin, continue
		c.Next()
	}
}

// SuperAdminAPIRequired middleware for API endpoints that require super admin privileges
func SuperAdminAPIRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		authUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Check if user has super admin privileges
		if authUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Super admin privileges required",
			})
			c.Abort()
			return
		}

		// User is super admin, continue
		c.Next()
	}
}

// RoleRequired middleware checks if the authenticated user has any of the specified roles
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		authUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user data",
			})
			c.Abort()
			return
		}

		// Check if user has any of the allowed roles
		hasRole := false
		for _, role := range allowedRoles {
			if authUser.Role == role || (role == "admin" && authUser.IsAdmin) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient privileges",
			})
			c.Abort()
			return
		}

		// User has required role, continue
		c.Next()
	}
}
