package handlers

import (
	"fmt"
	"net/http"
	"securetask/database"
	"securetask/models"

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
	database.DB.Where("user_id = ?", userID).Preload("User").Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

// VULNERABILITY #2: No input validation or sanitization
func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	// VULNERABILITY #3: No sanitization - XSS possible through description
	task := models.Task{
		Title:       req.Title,       // No sanitization
		Description: req.Description, // XSS vulnerability!
		Priority:    req.Priority,
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
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VULNERABILITY #3: No sanitization on updated fields
	database.DB.Model(&task).Updates(updates)

	c.JSON(http.StatusOK, task)
}

// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
// VULNERABILITY #2: No authorization check
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	// Anyone can delete any task!
	if err := database.DB.Delete(&models.Task{}, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// VULNERABILITY #1: SQL Injection in search functionality
// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
func SearchTasks(c *gin.Context) {
	searchTerm := c.Query("q")

	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term required"})
		return
	}

	// VULNERABILITY #1: Direct string concatenation - SQL INJECTION!
	query := fmt.Sprintf("SELECT * FROM tasks WHERE title LIKE '%%%s%%' OR description LIKE '%%%s%%'", searchTerm, searchTerm)

	// Using raw SQL without parameterization
	results, err := database.ExecuteRawSQL(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
