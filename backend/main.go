package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"securetask/database"
	"securetask/handlers"
	"securetask/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 4)
	return func(c *gin.Context) {

		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Limite exceed",
			})
		}

	}
}

func main() {
	// Load environment variables
	godotenv.Load()
	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

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
	r.Use(RateLimiter())

	// VULNERABILITY: Permissive CORS - allows all origins
	/*
	   Fix: Restrict CORS to only allow trusted origins and specific methods/headers.
	   How: Update the cors.Config to specify allowed origins, methods, and headers instead of allowing all.
	*/
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes (no authentication required)
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	// Protected routes (with auth middleware)
	authorized := r.Group("/api")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.GET("/tasks", handlers.GetTasks)
		authorized.POST("/tasks", handlers.CreateTask)
		authorized.PUT("/tasks/:id", handlers.UpdateTask)
		authorized.GET("/users/me", handlers.GetCurrentUser)

		authorized.GET("/api/tasks/search", handlers.SearchTasks)
		authorized.DELETE("/api/tasks/:id", handlers.DeleteTask)
		authorized.PUT("/api/users/:id/profile", handlers.UpdateProfile)

		admin := authorized.Group("/admin", handlers.AdminMiddleware())
		{
			admin.GET("/users", handlers.GetAllUsers)
		}
	}

	log.Println("🚀 Server starting on port 8080...")
	log.Println("⚠️  WARNING: This server contains intentional security vulnerabilities!")
	r.Run(":8080")
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return string(bytes)
}

func seedData() {
	seedUsers := []models.User{
		{
			Email:    "admin@example.com",
			Password: HashPassword("admin123"),
			Name:     "Admin User",
			Role:     "admin",
			Bio:      "I'm the administrator",
		},
		{
			Email:    "user@example.com",
			Password: HashPassword("password123"),
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
