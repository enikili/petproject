package taskService

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(id uint, task Task) (Task, error)
	GetUserTasks(id uint) ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
	GetAllTasks() ([]Task, error)
	GetTaskByID(id uint) (*Task, error)
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

func (r *taskRepository) GetUserTasks(id uint) ([]Task, error) {
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
	if task.Task != nil {
		existingTask.Task = task.Task
	}
	if *task.IsDone {
		existingTask.IsDone = task.IsDone
	}
	if err := r.db.Save(&existingTask).Error; err != nil {
		return Task{}, err
	}
	return existingTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	result := r.db.Delete(&Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTaskByID(id uint) (*Task, error) {
	var task Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}