package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"sso-web-app/internal/handlers"
	"sso-web-app/internal/middleware"
	"sso-web-app/internal/services"
)

func main() {
	// Load configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize services
	authService := services.NewAuthService()
	oauthService := services.NewOAuthService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, oauthService)
	adminHandler := handlers.NewAdminHandler()

	// Setup Gin router
	router := gin.Default()

	// Load HTML templates with proper parsing
	router.LoadHTMLGlob("templates/**/*.html")

	// Serve static files
	router.Static("/static", "./static")

	// Public routes
	public := router.Group("/")
	{
		public.GET("/", authHandler.Home)
		public.GET("/login", authHandler.LoginPage)
		public.POST("/login", authHandler.Login)
		public.GET("/register", authHandler.RegisterPage)
		public.POST("/register", authHandler.Register)
		public.GET("/logout", authHandler.Logout)
		
		// OAuth routes
		public.GET("/auth/google", authHandler.GoogleLogin)
		public.GET("/auth/google/callback", authHandler.GoogleCallback)
		public.GET("/auth/github", authHandler.GitHubLogin)
		public.GET("/auth/github/callback", authHandler.GitHubCallback)
	}

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/dashboard", authHandler.Dashboard)
		protected.GET("/profile", authHandler.Profile)
		protected.POST("/profile", authHandler.UpdateProfile)
	}

	// API routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/user", authHandler.GetUser)
		api.PUT("/user", authHandler.UpdateUser)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminRequired())
	{
		admin.GET("/dashboard", adminHandler.Dashboard)
		admin.GET("/users", adminHandler.UsersList)
		admin.GET("/users/:id", adminHandler.UserDetail)
	}

	// Admin API routes
	adminAPI := router.Group("/admin/api")
	adminAPI.Use(middleware.AuthMiddleware(), middleware.AdminAPIRequired())
	{
		adminAPI.PUT("/users/:id", adminHandler.UpdateUser)
		adminAPI.POST("/users/:id/activate", adminHandler.ActivateUser)
		adminAPI.POST("/users/:id/deactivate", adminHandler.DeactivateUser)
		adminAPI.DELETE("/users/:id", adminHandler.DeleteUser)
		adminAPI.POST("/users/:id/promote", adminHandler.PromoteToAdmin)
		adminAPI.POST("/users/:id/demote", adminHandler.DemoteFromAdmin)
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
