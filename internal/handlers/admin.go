package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sso-web-app/internal/models"
	"sso-web-app/internal/services"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		adminService: services.NewAdminService(),
	}
}

// Dashboard displays admin dashboard with statistics
func (h *AdminHandler) Dashboard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"title": "Unauthorized",
			"error": "Authentication required",
		})
		return
	}

	adminUser := user.(*models.User)
	
	stats, err := h.adminService.GetUserStats(adminUser)
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"title": "Access Denied",
				"error": "Admin privileges required",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Error",
			"error": "Failed to load dashboard data",
		})
		return
	}

	c.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"title":     "Admin Dashboard",
		"user":      adminUser,
		"stats":     stats,
		"isAdmin":   true,
		"activePage": "dashboard",
	})
}

// UsersList displays paginated list of all users
func (h *AdminHandler) UsersList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"title": "Unauthorized",
			"error": "Authentication required",
		})
		return
	}

	adminUser := user.(*models.User)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	// Parse filter parameters
	role := c.Query("role")
	search := c.Query("search")

	var users []*models.User
	var err error

	if search != "" {
		users, err = h.adminService.SearchUsers(adminUser, search, limit, offset)
	} else if role != "" {
		users, err = h.adminService.GetUsersByRole(adminUser, role, limit, offset)
	} else {
		users, err = h.adminService.GetAllUsers(adminUser, limit, offset)
	}

	if err != nil {
		if err == services.ErrNotAuthorized {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"title": "Access Denied",
				"error": "Admin privileges required",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Error",
			"error": "Failed to load users",
		})
		return
	}

	c.HTML(http.StatusOK, "admin-users.html", gin.H{
		"title":      "User Management",
		"user":       adminUser,
		"users":      users,
		"isAdmin":    true,
		"activePage": "users",
		"currentPage": page,
		"searchQuery": search,
		"roleFilter":  role,
	})
}

// UserDetail displays detailed view of a specific user
func (h *AdminHandler) UserDetail(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"title": "Unauthorized",
			"error": "Authentication required",
		})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Error",
			"error": "Invalid user ID",
		})
		return
	}

	targetUser, err := h.adminService.GetUserByID(adminUser, uint(userID))
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"title": "Access Denied",
				"error": "Admin privileges required",
			})
			return
		}
		if err == services.ErrUserNotFound {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"title": "User Not Found",
				"error": "The requested user was not found",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "Error",
			"error": "Failed to load user data",
		})
		return
	}

	c.HTML(http.StatusOK, "admin-user-detail.html", gin.H{
		"title":      "User Details",
		"user":       adminUser,
		"targetUser": targetUser,
		"isAdmin":    true,
		"activePage": "users",
	})
}

// UpdateUser handles user updates from admin
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	updatedUser, err := h.adminService.UpdateUser(adminUser, uint(userID), req)
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err == services.ErrInvalidRole {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser.ToResponse(),
	})
}

// DeactivateUser deactivates a user account
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updatedUser, err := h.adminService.DeactivateUser(adminUser, uint(userID))
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deactivated successfully",
		"user":    updatedUser.ToResponse(),
	})
}

// ActivateUser activates a user account
func (h *AdminHandler) ActivateUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updatedUser, err := h.adminService.ActivateUser(adminUser, uint(userID))
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User activated successfully",
		"user":    updatedUser.ToResponse(),
	})
}

// DeleteUser permanently deletes a user account
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.adminService.DeleteUser(adminUser, uint(userID))
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// PromoteToAdmin promotes a user to admin role
func (h *AdminHandler) PromoteToAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updatedUser, err := h.adminService.PromoteToAdmin(adminUser, uint(userID))
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user":    updatedUser.ToResponse(),
	})
}

// DemoteFromAdmin removes admin privileges from a user
func (h *AdminHandler) DemoteFromAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	adminUser := user.(*models.User)
	
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updatedUser, err := h.adminService.DemoteFromAdmin(adminUser, uint(userID))
	if err != nil {
		if err == services.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			return
		}
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin privileges removed successfully",
		"user":    updatedUser.ToResponse(),
	})
}
