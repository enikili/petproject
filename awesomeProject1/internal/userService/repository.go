package userService

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(id uint, user User) (User, error)
	DeleteUser(id uint) error
	GetUserByID(id uint) (*User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return User{}, ErrInvalidInput
	}

	result := r.db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUser(id uint, user User) (User, error) {
	if id == 0 {
		return User{}, ErrInvalidInput
	}

	var existingUser User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	// Обновляем только заполненные поля
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}

	if err := r.db.Save(&existingUser).Error; err != nil {
		return User{}, err
	}

	return existingUser, nil
}

func (r *userRepository) DeleteUser(id uint) error {
	if id == 0 {
		return ErrInvalidInput
	}

	result := r.db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) GetUserByID(id uint) (*User, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}

	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}