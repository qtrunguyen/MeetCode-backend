package service

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserService handles user-related operations
type UserService interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
}

type userService struct {
	db *sql.DB
}

// NewUserService creates a new instance of UserService
func NewUserService(db *sql.DB) UserService {
	return &userService{db: db}
}

// User model definition
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUp handles user sign-up
func (us *userService) SignUp(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Insert user into the database
	_, err := us.db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User signed up successfully"})
}

func (us *userService) LogIn(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Query the database for the user
	var storedPassword string
	err := us.db.QueryRow("SELECT password FROM users WHERE username = $1", credentials.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}

	// Verify the password
	if storedPassword != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
