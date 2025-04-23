package taskService



import (
	"errors"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	repo TaskRepository
}

func NewService(repo *taskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(id uint, task Task) (Task, error) {
	return s.repo.CreateTask(id, task)
}

func (s *TaskService) GetTaskByUserId(id uint) ([]Task, error) {
	return s.repo.GetTaskByUserId(id)
}

func (s *TaskService) UpdateTaskByID(id uint, task Task) (Task, error) {
	return s.repo.UpdateTaskByID(id, task)
}

func (s *TaskService) DeleteTask(id uint) error {

	return s.repo.DeleteTaskByID(id)
}