package taskService

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(userID uint, task Task) (Task, error) {
	if task.Task == nil || *task.Task == "" {
		return Task{}, errors.New("task description cannot be empty")
	}
	if userID == 0 {
		return Task{}, errors.New("user ID cannot be zero")
	}

	return s.repo.CreateTask(userID, task)
}

func (s *TaskService) GetTasksId(userID uint) ([]Task, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be zero")
	}

	tasks, err := s.repo.GetUserTasks(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) UpdateTaskByID(taskID uint, updates Task) (Task, error) {
	if taskID == 0 {
		return Task{}, errors.New("task ID cannot be zero")
	}
	updatedTask, err := s.repo.UpdateTaskByID(taskID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, err
	}
	return updatedTask, nil
}

func (s *TaskService) DeleteTasksById(taskID uint) error {
	if taskID == 0 {
		return errors.New("task ID cannot be zero")
	}

	err := s.repo.DeleteTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	tasks, err := s.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetTaskByID(taskID uint) (*Task, error) {
	if taskID == 0 {
		return nil, errors.New("task ID cannot be zero")
	}

	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}