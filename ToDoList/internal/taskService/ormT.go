package taskService

import "gorm.io/gorm"

type Task struct {
    gorm.Model
    Task   *string `json:"task"`
    IsDone *bool   `json:"is_done"`
    UserID *uint   `json:"user_id"`
}


