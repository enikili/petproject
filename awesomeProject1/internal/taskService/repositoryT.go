package taskService

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(id uint, task Task) (Task, error)
	GetTaskByUserId(id uint) ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(id uint, task Task) (Task, error) {
	task.UserID = &id
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetTaskByUserId(id uint) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("user_id = ?", id).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	var existingTask Task
	if err := r.db.First(&existingTask, id).Error; err != nil {
		return Task{}, err
	}
	existingTask.Task = task.Task
	existingTask.IsDone = task.IsDone

	if err := r.db.Save(&existingTask).Error; err != nil {
		return Task{}, err
	}

	return existingTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	var task Task
	result := r.db.Where("id = ?", id).Delete(&task)
	if result.Error != nil {
		return result.Error
	}
	return nil

}