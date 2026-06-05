package handlers

import (
	"net/http"
	"os"
	"securetask/database"
	"securetask/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Fix: read JWT secret lazily (at use time), not at package load time.
// godotenv.Load() runs inside main(), which executes AFTER this package's
// init. Reading os.Getenv here at load time would always be empty.
func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	return []byte(secret)
}

const authCookieName = "access_token"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// VULNERABILITY #2: No input validation, no password strength requirements
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isStrongPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password does not meet strength requirements"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashed),
		Name:     req.Name,
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// VULNERABILITY #2: No rate limiting - vulnerable to brute force attacks
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		authCookieName,
		tokenString,
		60*60*24, // 1 day
		"/",
		"",
		false,
		true,
	)

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(authCookieName)

		if err != nil || tokenString == "" {
			tokenString = c.GetHeader("Authorization")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
				c.Abort()
				return
			}

			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}
		}

		token, err := decodeJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if uid, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(uid))
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}

		c.Next()
	}

}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleI, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		roleStr, ok := roleI.(string)
		if !ok || roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func isStrongPassword(p string) bool {
	var hasMin = len(p) >= 8
	var hasUpper, hasLower, hasNumber bool
	for _, r := range p {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasNumber = true
		}
	}
	return hasMin && hasUpper && hasLower && hasNumber
}

func verifySignature(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret(), nil
	})
}

func decodeJWT(tokenString string) (*jwt.Token, error) {
	return verifySignature(tokenString)
}
