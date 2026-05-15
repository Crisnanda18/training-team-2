package main

import (
	"log"
	"securetask/database"
	"securetask/handlers"
	"securetask/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// VULNERABILITY #4: Hardcoded credentials directly in source code
const (
	// These should NEVER be hardcoded in production!
	DB_CONNECTION = "host=localhost user=taskuser password=taskpass123 dbname=securetask port=5432 sslmode=disable"
	JWT_SECRET    = "supersecret123"  // VULNERABILITY: Weak, hardcoded JWT secret
	ADMIN_KEY     = "admin-key-12345" // VULNERABILITY: Hardcoded API key
)

func main() {
	// Initialize database
	database.Connect()

	// Auto-migrate models
	err := database.DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data
	seedData()

	// Setup Gin router
	r := gin.Default()

	// VULNERABILITY: Permissive CORS - allows all origins
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes (no authentication required)
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	// VULNERABILITY #2: No authentication middleware on these routes!
	r.GET("/api/tasks/search", handlers.SearchTasks)        // Should require auth
	r.DELETE("/api/tasks/:id", handlers.DeleteTask)         // Should require auth
	r.PUT("/api/users/:id/profile", handlers.UpdateProfile) // Should require auth

	// VULNERABILITY #2: Admin route with no authorization check
	r.GET("/api/admin/users", handlers.GetAllUsers) // Anyone can access!

	// Protected routes (with auth middleware)
	authorized := r.Group("/api")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.GET("/tasks", handlers.GetTasks)
		authorized.POST("/tasks", handlers.CreateTask)
		authorized.PUT("/tasks/:id", handlers.UpdateTask)
		authorized.GET("/users/me", handlers.GetCurrentUser)
	}

	log.Println("🚀 Server starting on port 8080...")
	log.Println("⚠️  WARNING: This server contains intentional security vulnerabilities!")
	r.Run(":8080")
}

func seedData() {
	// Check if users already exist
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	// VULNERABILITY #5: Passwords stored in plain text (no hashing!)
	users := []models.User{
		{
			Email:    "admin@example.com",
			Password: "admin123", // Plain text password!
			Name:     "Admin User",
			Role:     "admin",
			Bio:      "I'm the administrator",
		},
		{
			Email:    "user@example.com",
			Password: "password123", // Plain text password!
			Name:     "Regular User",
			Role:     "user",
			Bio:      "Just a regular user",
		},
	}

	for _, user := range users {
		database.DB.Create(&user)
	}

	log.Println("✅ Database seeded with initial users")
}
