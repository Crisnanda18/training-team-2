package handlers

import (
	"html"
	"net/http"
	"securetask/database"
	"securetask/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

func GetTasks(c *gin.Context) {
	userID := c.GetUint("user_id")

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

// VULNERABILITY #2: No input validation or sanitization
// VULNERABILITY #3: No sanitization - XSS possible through description
/*
   Fix: Validate required fields and sanitize user-controlled strings before saving.
   How: Use `ShouldBindJSON` for validation, `strings.TrimSpace` to normalize input,
   and `html.EscapeString` to neutralize HTML/JS payloads that could trigger XSS when
   rendered in the frontend.
*/
func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	// sanitize inputs
	title := strings.TrimSpace(html.EscapeString(req.Title))
	description := html.EscapeString(req.Description)

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	task := models.Task{
		Title:       title,
		Description: description,
		Priority:    strings.TrimSpace(req.Priority),
		Status:      "todo",
		UserID:      userID,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// VULNERABILITY #2: No input validation
	// VULNERABILITY #3: No sanitization on updated fields
	/*
	   Fix: Restrict which fields can be updated, validate/sanitize string inputs,
	   and prevent updating sensitive fields like `user_id` or `id`.
	   How: Accept a JSON object, filter by an allow-list, escape string values,
	   then perform the update via GORM.
	*/
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allowed := map[string]bool{"title": true, "description": true, "priority": true, "status": true}
	safeUpdates := make(map[string]interface{})

	for k, v := range updates {
		key := strings.ToLower(k)
		if !allowed[key] {
			continue
		}

		// sanitize string values
		if s, ok := v.(string); ok {
			s = strings.TrimSpace(html.EscapeString(s))
			if key == "title" && s == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
				return
			}
			safeUpdates[key] = s
		} else {
			safeUpdates[key] = v
		}
	}

	if len(safeUpdates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields to update"})
		return
	}

	if err := database.DB.Model(&task).Updates(safeUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// reload task to return current state
	database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task)

	c.JSON(http.StatusOK, task)
}

// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
// VULNERABILITY #2: No authorization check
/*
   Fix: Require authentication and ensure only the owner can delete their task.
   How: Use `user_id` from the context, fetch the task filtered by `id` and `user_id`,
   and then delete if it belongs to the requester.
*/
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or not owned by user"})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// VULNERABILITY #1: SQL Injection in search functionality
// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
/*
   Fix: Use parameterized queries instead of string concatenation when executing raw SQL.
   How: Build the SQL with placeholders and pass the sanitized/wrapped search term as an argument.
*/
func SearchTasks(c *gin.Context) {
	searchTerm := c.Query("q")

	if strings.TrimSpace(searchTerm) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term required"})
		return
	}

	// parameterized query to prevent SQL injection
	query := "SELECT * FROM tasks WHERE title ILIKE ? OR description ILIKE ?"
	like := "%" + searchTerm + "%"

	results, err := database.ExecuteRawSQL(query, like, like)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
