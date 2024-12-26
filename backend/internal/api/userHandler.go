package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// SignUp is a handler for POST /signup
func Register(c *gin.Context) {
	// get the email/password off the request body
	var body *models.SignUpInput

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid json request format",
		})
		return
	}

	// check if email and password are provided
	if body.Email == "" || body.Password == "" || body.Name == "" || body.PasswordConfirm == "" {
		c.JSON(400, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	if body.Password != body.PasswordConfirm {
		c.JSON(400, gin.H{
			"error": "password do not match",
		})
		return
	}

	// check if the email is valid
	_, err = VerifyEmailExistence(body.Email)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return
	}

	// create a new user
	now := time.Now()
	user := models.User{
		Name: body.Name, 
		Email: strings.ToLower(body.Email), 
		Password: string(password), 
		Verified: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already exists",
		})
		return
	}
	// respond with the user
	c.JSON(200, gin.H{
		"status": "success",
	})
}

// Login is a handler for POST /login
func Login(c *gin.Context) {
	// get the email/password off the request body
	var body *models.SignInInput

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid json request format",
		})
		return
	}

	// Check if email and password are provided
	if body.Email == "" || body.Password == "" {
		c.JSON(400, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	// find the user
	var user models.User
	initializers.DB.First(&user, "email = ?", strings.ToLower(body.Email))
	if user.ID == [16]byte{} {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// generate a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// set the token in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*24*30, "/", "", false, true)
	c.JSON(200, gin.H{
		"status": "success",
	})
}

// Logout is a handler for GET /logout
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
}

// GetUser is a handler for GET /user
func GetUser(c *gin.Context) {
	user, _ := c.MustGet("userData").(models.User)
	userResponse := &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:  	   user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	//TODO: complete the implementation with user model
	c.JSON(200, gin.H{
		"status": "success",
		"user": userResponse,
	})
}
