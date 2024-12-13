package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// SignUp is a handler for POST /signup
func SignUp(c *gin.Context) {
	// get the email/password off the request body
	var body struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"pwd"`
	}

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid json request format",
		})
		return
	}

	// check if email and password are provided
	if body.Email == "" || body.Password == "" || body.FirstName == "" || body.LastName == "" {
		c.JSON(400, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	// hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return
	}

	// create a new user
	user := models.User{FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Password: string(password)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	// respond with the user
	c.JSON(200, gin.H{
		"status": "success",
		"user":   user,
	})
}

// Login is a handler for POST /login
func Login(c *gin.Context) {
	// get the email/password off the request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"pwd"`
	}

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
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
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
		"user":   user,
	})
}

// Logout is a handler for GET /logout
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
}

// GetUser is a handler for GET /user
func GetUser(c *gin.Context) {
	user, _ := c.Get("user")
	//TODO: complete the implementation with user model
	c.JSON(200, user)
}
