package userService

import (
	"awesomeProject1/internal/taskService"

	"gorm.io/gorm"
)
type Email string
type User struct {
	gorm.Model
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Tasks    []taskService.Task `json:"tasks" gorm:"foreignKey:UserID"`
}
