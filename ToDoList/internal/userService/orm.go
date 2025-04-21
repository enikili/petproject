package userService

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name     string `gorm:"not null" json:"name"`
    Email    string `gorm:"unique;not null" json:"email"`
    Password string `gorm:"not null" json:"-"`
}

