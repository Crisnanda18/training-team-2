package handlers

import (
	"net/http"
	"securetask/database"
	"securetask/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// VULNERABILITY #2: Returning password in response
	c.JSON(http.StatusOK, user)
}

// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
// VULNERABILITY #2: No authorization check - can update any user's profile
func UpdateProfile(c *gin.Context) {
	userID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VULNERABILITY #2: Anyone can update anyone's profile!
	// No check if the authenticated user matches the profile being updated

	// VULNERABILITY #3: Bio field not sanitized - XSS vulnerability
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Model(&user).Updates(updates)

	c.JSON(http.StatusOK, user)
}

// VULNERABILITY #2: No authentication required (exposed publicly in main.go)
// VULNERABILITY #2: No authorization check - anyone can access admin endpoint
func GetAllUsers(c *gin.Context) {
	// Should check if user has admin role, but doesn't!

	var users []models.User
	database.DB.Find(&users)

	// VULNERABILITY #2: Returning all users with passwords!
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}
