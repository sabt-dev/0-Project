package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID `gorm:"type:VARCHAR(36);primary_key"`
	Name             string    `gorm:"type:varchar(255);not null"`
	Email     		 string    `gorm:"unique;not null"`
	VerificationCode string
	Verified         bool      `gorm:"not null"`
	Password  		 string    `gorm:"not null" validate:"required,min=8"`
	Role             string    `gorm:"not null;type:varchar(255);default:'user'"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
}

// BeforeCreate hook will be called before inserting a new record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
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
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}