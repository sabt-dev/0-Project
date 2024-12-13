package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" validate:"required,min=2"`
	LastName  string `gorm:"not null" validate:"required,min=2"`
	Email     string `gorm:"unique;not null" validate:"required,email,min=8"`
	Password  string `json:"-" gorm:"not null" validate:"required,min=8"`
}
