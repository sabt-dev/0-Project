package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/models"
	"github.com/sabt-dev/0-Project/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// SignUp is a handler for POST /signup
func Register(c *gin.Context) {
	// get the email/password off the request body
	var body *models.SignUpInput
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"error": "Invalid json request format",
		})
		return
	}

	// check if email and password are provided
	if body.Email == "" || body.Password == "" || body.Name == "" || body.PasswordConfirm == "" {
		c.JSON(400, gin.H{
			"status": "fail",
			"error": "Email and password are required",
		})
		return
	}

	if body.Password != body.PasswordConfirm {
		c.JSON(400, gin.H{
			"status": "fail",
			"error": "password do not match",
		})
		return
	}

	// check if the domain of email and the format are valid
	_, err = utils.VerifyEmailExistence(body.Email)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
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
	var now time.Time = time.Now()
	newUser := models.User{
		ID: uuid.New(),
		Name: body.Name, 
		Email: strings.ToLower(body.Email), 
		Password: string(password), 
		Role: "user",
		Verified: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result := initializers.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": "user already exists or failed to create user",
		})
		return
	}

	// Generate Verification Code & Update User in Database
	var verCode string = utils.GenerateEmailVerificationCode()
	var encodedVerCode string = utils.Encode(verCode)

	newUser.VerificationCode = encodedVerCode
	initializers.DB.Save(newUser)

	var firstName string = newUser.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// Send Verification Code to User's Email
	err = utils.SendVerificationCode(newUser.Email, verCode, firstName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Failed to send verification code",
		})
		return
	}

	// respond with the user
	c.JSON(201, gin.H{
		"status": "success",
		"message": "verification code sent to your email " + newUser.Email,
		"uid": newUser.ID,
		"email": newUser.Email,
	})
}

// Login is a handler for POST /login
func Login(c *gin.Context) {
	// get the email/password off the request body
	var body *models.SignInInput

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": "Invalid json request format",
		})
		return
	}

	// Check if email and password are provided
	if body.Email == "" || body.Password == "" {
		c.JSON(400, gin.H{
			"status": "fail",
			"error": "Email and password are required",
		})
		return
	}

	// find the user
	var user models.User
	initializers.DB.First(&user, "email = ?", strings.ToLower(body.Email))
	if user.ID == [16]byte{} {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": "Invalid email or password",
		})
		return
	}

	// check if the user is verified
	if !user.Verified {
		c.JSON(http.StatusConflict, gin.H{
			"status": "fail",
			"error": "Invalid email or password",
		})
		return
	}

	// compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error": "Invalid email or password",
		})
		return
	}

	// generate a token
	tokenString, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Something went wrong",
		})
		return
	}

	// set the token in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*24*30, "/", "", false, true)
	c.JSON(200, gin.H{
		"status": "success",
		"message": "Logged in successfully",
	})
}

// VerifyEmail is a handler for GET /verifyemail?code=[verification_code]
func VerifyUserEmail(c *gin.Context) {
	var code string = c.Query("code")
	var verification_code string = utils.Encode(code)
	var updatedUser models.User
	result := initializers.DB.First(&updatedUser, "verification_code = ?", verification_code)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", 
			"message": "Invalid verification code or user doesn't exists",
		})
		return
	}

	if updatedUser.Verified {
		c.JSON(http.StatusConflict, gin.H{
			"status": "fail", 
			"message": "User already verified",
		})
		return
	}

	var now time.Time = time.Now()
	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	updatedUser.VerifiedAt = &now
	initializers.DB.Save(&updatedUser)
	c.JSON(200, gin.H{
		"status": "success",
		"message": "User verified successfully, now you can login",
	})
}

// Logout is a handler for GET /logout
func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
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
		"data": userResponse,
	})
}
