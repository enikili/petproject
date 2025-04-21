package userService

import (
	"awesomeProject1/internal/taskService"
	"awesomeProject1/internal/web/tasks"

	"gorm.io/gorm"
)



type UserRepository interface {
	CreateUser(user User) (User, error)

	GetAllUsers() ([]User, error)

	GetUserByEmail(Email string) (User, error)

	GetTasksForUser(userID uint) ([]tasks.Task, error)

	UpdateUserByID(id uint, user User) (User, error)

	DeleteUserByID(id uint) error
}

type userRepository struct {
	db       *gorm.DB
	taskRepo taskService.TaskRepository
}

func NewUserRepository(db *gorm.DB, taskRepo taskService.TaskRepository) UserRepository {
	return &userRepository{db: db, taskRepo: taskRepo}
}

func (repo *userRepository) CreateUser(user User) (User, error) {
	result := repo.db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (repo *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := repo.db.Find(&users).Error
	return users, err
}

func (repo *userRepository) GetUserByEmail(email string) (User, error) {
	var user User
	result := repo.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (repo *userRepository) UpdateUserByID(id uint, user User) (User, error) {
	var existingUser User

	if err := repo.db.First(&existingUser, id).Error; err != nil {
		return User{}, err
	}

	result := repo.db.Model(&existingUser).Updates(user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return existingUser, nil
}

func (repo *userRepository) DeleteUserByID(id uint) error {
	result := repo.db.Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *userRepository) GetTasksForUser(userID uint) ([]tasks.Task, error) {
	dbTasks, err := repo.taskRepo.GetTasksByUserID(userID)
	if err != nil {
		return nil, err
	}

	var Tasks []tasks.Task
	for _, dbTask := range dbTasks {
		Task := tasks.Task{
			Id:     dbTask.ID,
			Task:   *dbTask.Task,
			IsDone: dbTask.IsDone,
			UserId: Task.UserID,
		}
		Tasks = append(Tasks, Task)
	}

	return Tasks, nil
}