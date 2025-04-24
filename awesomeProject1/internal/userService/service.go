package userService

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input data")
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user User) (User, error) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return User{}, ErrInvalidInput
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) UpdateUser(id uint, updates User) (User, error) {
	if id == 0 {
		return User{}, ErrInvalidInput
	}

	// Проверяем, что есть хотя бы одно поле для обновления
	if updates.Username == "" && updates.Password == "" && updates.Email == "" {
		return User{}, ErrInvalidInput
	}

	return s.repo.UpdateUser(id, updates)
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return ErrInvalidInput
	}
	return s.repo.DeleteUser(id)
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.GetUserByID(id)
}