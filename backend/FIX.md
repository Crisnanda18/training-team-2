# List of Vulnerability

## SQL Injection

**Files to Review:**

- `backend/database/database.go`
  before:

```go
func ExecuteRawSQL(query string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// This allows SQL injection!
	rows, err := DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			entry[col] = values[i]
		}

		results = append(results, entry)
	}

	return results, nil
}
```

After:

```go
func ExecuteRawSQL(query string, args ...any) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := DB.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			entry[col] = values[i]
		}

		results = append(results, entry)
	}

	return results, nil
}
```

Explanation: The original code directly executes the raw SQL query without any parameterization, which can lead to SQL injection vulnerabilities. By modifying the function to accept additional arguments and using parameterized queries, we can mitigate this risk.

- `backend/handlers/tasks.go`

before:

```go
func SearchTasks(c *gin.Context) {
	searchTerm := c.Query("q")

	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term required"})
		return
	}

	query := fmt.Sprintf("SELECT * FROM tasks WHERE title LIKE '%%%s%%' OR description LIKE '%%%s%%'", searchTerm, searchTerm)

	results, err := database.ExecuteRawSQL(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
```

after:

```go
func SearchTasks(c *gin.Context) {
	searchTerm := c.Query("q")

	if strings.TrimSpace(searchTerm) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term required"})
		return
	}

	query := "SELECT * FROM tasks WHERE title ILIKE ? OR description ILIKE ?"
	like := "%" + searchTerm + "%"

	results, err := database.ExecuteRawSQL(query, like, like)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
```

## Authentication and Authorization

**Files to Review:**

- `backend/main.go` (route definitions)
  before:

```go
func main() {
	database.Connect()

	err := database.DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedData()

	r := gin.Default()

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
	r.GET("/api/tasks/search", handlers.SearchTasks)
	r.DELETE("/api/tasks/:id", handlers.DeleteTask)
	r.PUT("/api/users/:id/profile", handlers.UpdateProfile)

	r.GET("/api/admin/users", handlers.GetAllUsers)

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
```

after:

```go
func main() {
	godotenv.Load()

	database.Connect()

	err := database.DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedData()

	r := gin.Default()
	r.Use(RateLimiter())

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
```

- `backend/handlers/auth.go`
  before:

```go

```

after:

```go

```

- `backend/handlers/users.go`
