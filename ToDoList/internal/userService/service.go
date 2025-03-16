package userService

import (
 "errors"
)


type CreateUser struct {
 Name  string `json:"name"`
 Email string `json:"email"`
}


type UpdateUser struct {
 Name  *string `json:"name,omitempty"` 
 Email *string `json:"email,omitempty"`
}

type UserService struct {
 repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
 return &UserService{repo: repo}
}

func (s *UserService) CreateUser(userCreate *UserCreate) (*User, error) {
 if userCreate.Name == "" || userCreate.Email == "" {
  return nil, errors.New("name and email are required")
 }

 
 user := &User{
  Name:  userCreate.Name,
  Email: userCreate.Email,
 }

 createdUser, err := s.repo.CreateUser(user)
 if err != nil {
  return nil, err
 }

 
 userDTO := &User{
  ID:        createdUser.ID,
  Name:      createdUser.Name,
  Email:     createdUser.Email,
  CreatedAt: createdUser.CreatedAt,
  UpdatedAt: createdUser.UpdatedAt,
 }

 return userDTO, nil
}

func (s *UserService) GetAllUsers() ([]*User, error) {
 users, err := s.repo.GetAllUsers()
 if err != nil {
  return nil, err
 }

 
 userDTOs := make([]*User, len(users))
 for i, user := range users {
  userDTOs[i] = &User{
   ID:        user.ID,
   Name:      user.Name,
   Email:     user.Email,
   CreatedAt: user.CreatedAt,
   UpdatedAt: user.UpdatedAt,
  }
 }

 return userDTOs, nil
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
 user, err := s.repo.GetUserByID(id)
 if err != nil {
  return nil, err
 }

 
 userDTO := &User{
  ID:        user.ID,
  Name:      user.Name,
  Email:     user.Email,
  CreatedAt: user.CreatedAt,
  UpdatedAt: user.UpdatedAt,
 }

 return userDTO, nil
}

func (s *UserService) UpdateUser(id uint, userUpdate *UserUpdate) (*User, error) {
 
 existingUser, err := s.repo.GetUserByID(id)
 if err != nil {
  return nil, err
 }

 
 if userUpdate.Name != nil {
  existingUser.Name = *userUpdate.Name
 }
 if userUpdate.Email != nil {
  existingUser.Email = *userUpdate.Email
 }

 updatedUser, err := s.repo.UpdateUserByID(id, existingUser)
 if err != nil {
  return nil, err
 }

 
 userDTO := &User{
  ID:        updatedUser.ID,
  Name:      updatedUser.Name,
  Email:      updatedUser.Email,
  CreatedAt: updatedUser.CreatedAt,
  UpdatedAt: updatedUser.UpdatedAt,
 }

 return userDTO, nil
}

func (s *UserService) DeleteUser(id uint) error {
 return s.repo.DeleteUserByID(id)
}