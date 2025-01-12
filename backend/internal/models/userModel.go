package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID `gorm:"type:VARCHAR(36);primary_key"`
	FirstName            string    `gorm:"type:varchar(255);not null"`
	LastName             string    `gorm:"type:varchar(255);not null"`
	Username			 string    `gorm:"unique;not null"`
	Email                string    `gorm:"unique;not null"`
	VerificationCode     *string
	Verified             bool 	   `gorm:"not null"`
	VerifiedAt           *time.Time
	Password             string    `gorm:"not null" validate:"required,min=8"`
	Role                 string    `gorm:"not null;type:varchar(255);default:'user'"`
	PasswordResetToken   *string   `gorm:"type:varchar(255)"`
	PasswordResetExpires *time.Time
	RefreshToken         *string    `gorm:"type:varchar(255)"`
    RefreshTokenExpiresAt *time.Time
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

type SignUpInput struct {
	FirstName       string `json:"fname" binding:"required"`
	LastName        string `json:"lname" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	FirstName string    `json:"fname,omitempty"`
	LastName  string    `json:"lname,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
