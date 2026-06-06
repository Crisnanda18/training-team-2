package handlers

import (
	"html"
	"net/http"
	"securetask/database"
	"securetask/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.UserResponse
	if err := database.DB.Model(&models.User{}).First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// VULNERABILITY #2: Returning password in response
	/*
		fix: Use a separate struct for responses that doesn't include password
	*/
	c.JSON(http.StatusOK, user)
}

// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
// VULNERABILITY #2: No authorization check - can update any user's profile
func UpdateProfile(c *gin.Context) {
	targetUserID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	authenticatedUserID := c.GetUint("user_id")
	role, _ := c.Get("role")
	if authenticatedUserID != uint(targetUserID) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	safeUpdates := map[string]interface{}{}
	if name, ok := updates["name"].(string); ok {
		cleanName := strings.TrimSpace(html.EscapeString(name))
		if cleanName != "" {
			safeUpdates["name"] = cleanName
		}
	}
	if bio, ok := updates["bio"].(string); ok {
		safeUpdates["bio"] = html.EscapeString(strings.TrimSpace(bio))
	}

	if len(safeUpdates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.DB.Model(&user).Updates(safeUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
// VULNERABILITY #2: No authorization check - anyone can access admin endpoint
func GetAllUsers(c *gin.Context) {
	var users []models.UserResponse
	if err := database.DB.Model(&models.User{}).
		Select("id", "email", "name", "role", "bio").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}
