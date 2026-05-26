package main

import (
	"log"
	"securetask/database"
	"securetask/handlers"
	"securetask/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

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
	/*
	   Fix: Seed data idempotently instead of returning early when the users table is not empty.
	   How: Upsert the known users by email, capture their actual IDs, then create tasks for those
	   IDs with error checks so seeding still works on a partially populated database.
	*/
	seedUsers := []models.User{
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

	createdUsers := make(map[string]models.User)
	for _, seedUser := range seedUsers {
		var existing models.User
		err := database.DB.Where("email = ?", seedUser.Email).First(&existing).Error
		if err != nil {
			database.DB.Create(&seedUser)
			if seedErr := database.DB.Where("email = ?", seedUser.Email).First(&existing).Error; seedErr != nil {
				log.Printf("failed to seed user %s: %v", seedUser.Email, seedErr)
				continue
			}
		}

		createdUsers[seedUser.Email] = existing
	}

	seedTasks := []models.Task{
		{
			Title:       "Admin Task",
			Description: "This is a task for the admin user",
			Status:      "todo",
			Priority:    "high",
			UserID:      createdUsers["admin@example.com"].ID,
		},
		{
			Title:       "User Task",
			Description: "This is a task for the regular user",
			Status:      "in_progress",
			Priority:    "medium",
			UserID:      createdUsers["user@example.com"].ID,
		},
	}

	for _, seedTask := range seedTasks {
		if seedTask.UserID == 0 {
			log.Printf("skipping task seed %q because its user was not created", seedTask.Title)
			continue
		}

		var existing models.Task
		err := database.DB.Where("title = ? AND user_id = ?", seedTask.Title, seedTask.UserID).First(&existing).Error
		if err != nil {
			if createErr := database.DB.Create(&seedTask).Error; createErr != nil {
				log.Printf("failed to seed task %q: %v", seedTask.Title, createErr)
			}
		}
	}

	log.Println("✅ Database seeded with initial users and tasks")
}
