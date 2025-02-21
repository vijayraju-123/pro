package handlers

import (
	"net/http"
	"task-management-system/Backend/models"
	"task-management-system/Backend/utils"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash password (use bcrypt or similar)
	user.Password = utils.HashPassword(user.Password)

	// Save user to database (implement DB logic here)
	// Example: config.DB.Exec("INSERT INTO users ...", user.Username, user.Email, user.Password)
	// ...

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Verify user credentials (implement DB logic here)
	// Example: rows, _ := config.DB.Query("SELECT * FROM users WHERE email = ?", user.Email)
	// ...

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}