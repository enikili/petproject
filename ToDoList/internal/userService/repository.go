package userService

import (
 "gorm.io/gorm"
 "time"
)

// Users model


type User struct {
 ID        int64    `json:"id"`
 Name      string    `json:"name"`
 Email     string    `json:"email"`
 Password  string    `json:"-"`
 CreatedAt time.Time `json:"created_at"`
 UpdatedAt time.Time `json:"updated_at"`
}
type UserCreate struct {
 Name  string `json:"name"`
 Email string `json:"email"`
}
type UserUpdate struct {
 Name  *string `json:"name,omitempty"`
 Email *string `json:"email,omitempty"`
}
type UserRepository interface {
 CreateUser(user *User) (*User, error) 
 GetAllUsers() ([]*User, error)           
 GetUserByID(id uint) (*User, error)      
 UpdateUserByID(id uint, user *User) (*User, error) 
 DeleteUserByID(id uint) error
 PatchUserByID(id uint, updates map[string]interface{}) (*User, error) 
}

type userRepository struct {
 db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
 return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *User) (*User, error) {
 result := r.db.Create(user)
 if result.Error != nil {
  return nil, result.Error 
 }
 return user, nil
}

func (r *userRepository) GetAllUsers() ([]*User, error) {
 var users []*User // Слайс указателей
 err := r.db.Find(&users).Error
 return users, err
}

func (r *userRepository) GetUserByID(id uint) (*User, error) {
    var user User 
    result := r.db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil 
}

func (r *userRepository) UpdateUserByID(id uint, user *User) (*User, error) {
    var existingUser User
    result := r.db.First(&existingUser, id)
    if result.Error != nil {
        return nil, result.Error
    }

    
    existingUser.Name = user.Name
    existingUser.Email = user.Email

    
    result = r.db.Save(&existingUser)
    if result.Error != nil {
        return nil, result.Error
    }

    return &existingUser, nil
}

func (r *userRepository) DeleteUserByID(id uint) error {
    var user User
    result := r.db.First(&user, id)
    if result.Error != nil {
        return result.Error
    }

    result = r.db.Delete(&user)
    return result.Error
}

func (r *userRepository) PatchUserByID(id uint, updates map[string]interface{}) (*User, error) {
    var user User
    result := r.db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }

    
    updates["updated_at"] = time.Now()

    
    result = r.db.Model(&user).Updates(updates)
    if result.Error != nil {
        return nil, result.Error
    }

    return &user, nil
}