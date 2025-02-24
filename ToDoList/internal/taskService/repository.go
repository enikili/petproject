package taskService
import (
	
	"gorm.io/gorm"
)


type TaskRepository interface {
	CreateTask(task Tasks) (Tasks, error)
	GetAllTasks() ([]Tasks, error)
	UpdateTaskByID(id uint, task Tasks) (Tasks, error)
	DeleteTaskByID(id uint) error
	PatchTaskByID(id uint, updates map[string]interface{}) (Tasks, error)
}

type taskRepository struct {
	db *gorm.DB
}
 
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Tasks) (Tasks, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Tasks{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Tasks, error) {
	var tasks []Tasks
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, task Tasks) (Tasks, error) {
	// First, find the existing task by ID
	var existingTask Tasks
	result := r.db.First(&existingTask, id)
	if result.Error != nil {
		return Tasks{}, result.Error
	}

	existingTask.Title = task.Title
	existingTask.Description = task.Description
	

	result = r.db.Save(&existingTask)
	if result.Error != nil {
		return Tasks{}, result.Error
	}

	return existingTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	result := r.db.Delete(&Tasks{}, id)
	return result.Error
}

func (r *taskRepository) PatchTaskByID(id uint, updates map[string]interface{}) (Tasks, error) {
	var task Tasks
	result := r.db.First(&task, id)
	if result.Error != nil {
		return Tasks{}, result.Error
	}

	result = r.db.Model(&task).Updates(updates)
	if result.Error != nil {
		return Tasks{}, result.Error
	}

	return task, nil
}
