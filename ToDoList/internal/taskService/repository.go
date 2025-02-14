package taskService

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTasks() ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
	PatchTaskByID(id uint, updates map[string]interface{}) (Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	// First, find the existing task by ID
	var existingTask Task
	result := r.db.First(&existingTask, id)
	if result.Error != nil {
		return Task{}, result.Error
	}

	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.Completed = task.Completed

	result = r.db.Save(&existingTask)
	if result.Error != nil {
		return Task{}, result.Error
	}

	return existingTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	result := r.db.Delete(&Task{}, id)
	return result.Error
}

func (r *taskRepository) PatchTaskByID(id uint, updates map[string]interface{}) (Task, error) {
	var task Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		return Task{}, result.Error
	}

	result = r.db.Model(&task).Updates(updates)
	if result.Error != nil {
		return Task{}, result.Error
	}

	return task, nil
}
