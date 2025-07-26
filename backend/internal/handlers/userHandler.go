package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/models"
	"github.com/sabt-dev/0-Project/internal/utils"
)

// SignUp is a handler for POST /signup
func Register(c *gin.Context) {
	// get the email/password off the request body
	var body *models.SignUpInput
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  "Invalid json request format",
		})
		return
	}

	// check if email and password are provided
	if body.Email == "" || body.Password == "" || body.FirstName == "" || body.LastName == "" || body.Username == "" {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  "Missing required fields",
		})
		return
	}

	// check if the domain of email and the format are valid
	_, err = utils.VerifyEmailExistence(body.Email)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  "invalid email or email does not exist",
		})
		return
	}

	// hash the password
	password, err := utils.HashPassword(body.Password)
	if err != nil {
		return
	}

	// Start a new transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// create a new user
	var now time.Time = time.Now()
	newUser := models.User{
		ID:        uuid.New(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Username:  body.Username,
		Email:     strings.ToLower(body.Email),
		Password:  password,
		Role:      "user",
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := tx.Create(&newUser)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "user already exists or failed to create user",
		})
		return
	}

	// Generate Verification Code & Update User in Database
	var verCode string = utils.GenerateEmailVerificationCode()
	var encodedVerCode string = utils.Encode(verCode)

	newUser.VerificationCode = &encodedVerCode
	if result := tx.Save(newUser); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// Send Verification Code to User's Email
	err = utils.SendVerificationCode(newUser.Email, verCode, newUser.FirstName)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Failed to send verification code",
		})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// respond with the user
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "verification code sent to your email " + newUser.Email,
		"uid":     newUser.ID,
		"email":   newUser.Email,
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
			"error":  "Invalid json request format",
		})
		return
	}

	// Check if email and password are provided
	if body.Email == "" || body.Password == "" {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  "Email and password are required",
		})
		return
	}

	// check if the domain of email and the format are valid
	_, err = utils.VerifyEmailExistence(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid email or email does not exist",
		})
		return
	}

	// find the user
	var user models.User
	initializers.DB.First(&user, "email = ?", strings.ToLower(body.Email))
	if user.ID == [16]byte{} {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid email or password",
		})
		return
	}

	// check if the user is verified
	if !user.Verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid email or password",
		})
		return
	}

	// compare the password
	err = utils.CheckPasswordHash(body.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid email or password",
		})
		return
	}

	// Generate a JWT access token
	accessTokenString, err := utils.GenerateAccessToken(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// Generate a JWT refresh token
    refreshTokenString, err := utils.GenerateRefreshToken(&user, c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "Something went wrong",
        })
        return
    }

	// Start a new transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// Update the user's refresh token
	var expiresAt time.Time = time.Now().Add(7 * 24 * time.Hour)
    user.RefreshToken = &refreshTokenString
	user.RefreshTokenExpiresAt = &expiresAt // Refresh token valid for 7 days
	if err := tx.Save(&user); err.Error != nil {
		tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "Failed to save refresh token",
        })
        return
    }

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// set the token in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", accessTokenString, 60*60*24, "/", "", false, true)
	c.SetCookie("RefreshToken", refreshTokenString, 60*60*24*7, "/", "", false, true)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Logged in successfully",
		"data":    gin.H{
			"uid": user.ID,
		},
	})
}

// VerifyEmail is a handler for POST /verifyemail
func VerifyUserEmail(c *gin.Context) {
	var body struct {
		Code string `json:"code" binding:"required"`
	}
	
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid json request format",
		})
		return
	}
	
	var verification_code string = utils.Encode(body.Code)

	// Start a new transaction
    tx := initializers.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "something went wrong",
		})
        return
    }

	var updatedUser models.User
	result := tx.First(&updatedUser, "verification_code = ?", verification_code)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid verification code or user doesn't exists",
		})
		return
	}

	if updatedUser.Verified {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "User already verified",
		})
		return
	}

	var now time.Time = time.Now()
	updatedUser.VerificationCode = nil
	updatedUser.Verified = true
	updatedUser.VerifiedAt = &now
	if tx.Save(&updatedUser).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Somthing went wrong",
		})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "User verified successfully, now you can login",
	})
}

// RequestPasswordReset is a handler for POST /request-password-reset
func RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid json request format",
		})
		return
	}

	// Check if the email is valid
	isValid, _ := utils.VerifyEmailExistence(request.Email)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid email",
		})
		return
	}

	// Start a new transaction
    tx := initializers.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Something went wrong",
		})
        return
    }

	// Check if the user exists
	var user models.User
	if tx.First(&user, "email = ?", request.Email).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "User not found",
		})
		return
	}

	// Generate a password reset token
	expirationTime := time.Now().Add(1 * time.Hour)
	resetToken := uuid.New().String()
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpires = &expirationTime // Token expires in 1 hour
	if tx.Save(&user).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Something went wrong",
		})
		return
	}

	// Encode the reset token and Send the reset token to the user's email
	err := utils.SendPasswordResetEmail(user.Email, utils.URLencode(resetToken))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Failed to send password reset email",
		})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password reset code sent to your email",
	})
}

// ResetPassword is a handler for PUT /reset-password
func ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid json request format",
		})
		return
	}


	// Start a new transaction
    tx := initializers.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Something went wrong",
		})
        return
    }

	// Decode the reset token
	tokenString, err := utils.URLdecode(request.Token)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid or time expired",
		})
		return
	}

	var user models.User
	result := tx.First(&user, "password_reset_token = ?", tokenString)
	if result.Error != nil || user.PasswordResetExpires.Before(time.Now()) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "Invalid or time expired",
		})
		return
	}

	// Check if the new password is the same as the old password
	if err := utils.CheckPasswordHash(request.NewPassword, user.Password); err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "New password cannot be the same as the old password",
		})
		return
	}

	// Update the user's password
	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Something went wrong",
		})
		return
	}

	// Update the user's password
	user.Password = hashedPassword
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil
	if result := tx.Save(&user); result.Error != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error": "Failed to update password",
		})
        return
    }

	// Commit the transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password reset successfully",
	})
}

// Logout is a handler for GET /logout
func Logout(c *gin.Context) {
	user, exists := c.Get("userData")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "fail",
            "error":  "User not authenticated",
        })
        return
    }

    // Invalidate the refresh token
    var userData models.User = user.(models.User)
    userData.RefreshToken = nil
    userData.RefreshTokenExpiresAt = nil

    tx := initializers.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "Something went wrong",
        })
        return
    }

    if err := tx.Save(&userData).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "Something went wrong",
        })
        return
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "Something went wrong",
        })
        return
    }

    // Clear the cookies
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", "", -1, "/", "", false, true)
    c.SetCookie("RefreshToken", "", -1, "/", "", false, true)

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Logged out successfully",
    })
}

// GetUser is a handler for GET /user
func GetUser(c *gin.Context) {
	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "User not authenticated",
		})
		return
	}
	user, ok := userData.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Failed to parse user data",
		})
		return
	}
	userResponse := &models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
	}
	//TODO: complete the implementation with user model
	c.JSON(200, gin.H{
		"status": "success",
		"data":   userResponse,
	})
}

// RefreshToken is a handler for POST /refresh-token
func RefreshToken(c *gin.Context) {
    refreshToken, err := c.Cookie("RefreshToken")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "fail",
            "error":  "Refresh token not provided",
        })
        return
    }

    // Validate the refresh token
    claims, err := utils.ValidateToken(refreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "fail",
            "error":  "Invalid refresh token",
        })
        return
    }

	// check if the user has the same IP-address
	//if claims["ip"] != c.ClientIP() {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"status": "fail",
	//		"error":  "Invalid refresh token",
	//	})
	//	return
	//}

    // Extract user ID from token claims
    userID, ok := claims["sub"].(string)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "fail",
            "error":  "Invalid refresh token",
        })
        return
    }

    var user models.User
    result := initializers.DB.Where("id = ? AND refresh_token = ?", userID, refreshToken).First(&user)
    if result.Error != nil || user.RefreshTokenExpiresAt.Before(time.Now()) {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "fail",
            "error":  "Invalid or expired refresh token",
        })
        return
    }

    // Generate a new access token
	accessToken, err := utils.GenerateAccessToken(&user, c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "Something went wrong",
        })
        return
    }

    // Set the new access token in a cookie
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", accessToken, 60*60*24, "/", "", false, true)

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Access token refreshed successfully",
    })
}
