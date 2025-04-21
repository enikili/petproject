package taskService

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)

	GetAllTasks() ([]Task, error)

	GetTaskByID(id uint) (Task, error)

	GetTasksByUserID(userID uint) ([]Task, error)

	UpdateTaskByID(id uint, task Task) (Task, error)

	DeleteTaskByID(id uint) error
}

func (r *taskRepository) GetTasksByUserID(userID uint) ([]Task, error) {
 var tasks []Task
 result := r.db.Where("user_id = ?", userID).Find(&tasks)
 if result.Error != nil {
  return []Task{}, result.Error
 }
 return tasks, nil
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
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

func (repo *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	var existingTask Task

	if err := repo.db.First(&existingTask, id).Error; err != nil {
		return Task{}, err
	}

	result := repo.db.Model(&existingTask).Updates(task)
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


func (r *taskRepository) GetTaskByID(id uint) (*Task, error) {
    var task Task
    if err := r.db.First(&task, id).Error; err != nil {
        return nil, err
    }
    return &task, nil
}
