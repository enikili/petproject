package taskService

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Task        string
	IsDone      bool
	Title       string
	Description string
	Completed   *time.Time
}
