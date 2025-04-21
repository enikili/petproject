package userService

import (
	"awesomeProject1/internal/web/tasks"
	"errors"
	

)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) GetAllUsers() ([]User, error) {
	return service.repo.GetAllUsers()
}

func (service *UserService) GetTasksForUser(userID uint) ([]tasks.Task, error) {
	return service.repo.GetTasksForUser(userID)
}

func (service *UserService) CreateUser(user User) (User, error) {

	if user.Name == "" {
		return User{}, errors.New("username cannot be empty")
	}
	
	if user.Password == "" {
		return User{}, errors.New("password cannot be empty")
	}

	existingUser, err := service.repo.GetUserByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return User{}, errors.New("user with this email already exists")
	}

	return service.repo.CreateUser(user)
}

func (service *UserService) UpdateUserByID(id uint, user User) (User, error) {
	return service.repo.UpdateUserByID(id, user)
}

func (service *UserService) DeleteUserByID(id uint) error {
	return service.repo.DeleteUserByID(id)
}