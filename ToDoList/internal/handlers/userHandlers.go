package handlers

import (
	"context"
	"errors"
	"awesomeProject1/internal/userService"
	"awesomeProject1/internal/web"
	
)

type Handlers struct {
	Service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *Handlers { return &Handlers{Service: service} }

func (h *Handlers) GetUsers(_ context.Context, _ User.GetUsersRequestObject) (User.GetUsersResponseObject, error) {
	allUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}
	response := User.GetUsers200JSONResponse{}

	for _, user := range allUsers {
		user := User.User{
			Id:    &user.ID,
			Email: &user.Email,
			Name:  &user.Name,
		}
		response = append(response, user)
	}

	return response, nil
}

func (h *Handlers) GetUsersIdTasks(_ context.Context, request User.GetUsersIdTasksRequestObject) (User.GetUsersIdTasksResponseObject, error) {
	userID := request.Id

	dbTasks, err := h.Service.GetTasksForUser(uint(userID))
	if err != nil {
		return nil, err
	}

	var responseTasks []User.Task
	for _, task := range dbTasks {
		responseTask := User.Task{
			Id:     task.Id,
			Task:   "",
			IsDone: task.IsDone,
			UserId: 0,
		}
		if task.Task != "" {
			responseTask.Task = *&task.Task
		}

		responseTasks = append(responseTasks, responseTask)
	}

	return User.GetUsersIdTasks200JSONResponse(responseTasks), nil
}

func (h *Handlers) PostUsers(_ context.Context, request User.PostUsersRequestObject) (User.PostUsersResponseObject, error) {
	userRequest := request.Body

	userToCreate := userService.Users{

		Email:    *userRequest.Email,
		Name:     *userRequest.Name,
		Password: *userRequest.Password,
	}
	createdUser, err := h.Service.CreateUser(userToCreate)

	if userRequest.Password == nil {
		return nil, errors.New("password is required")
	}

	if err != nil {
		return nil, err
	}

	response := User.PostUsers201JSONResponse{
		Id:    &createdUser.ID,
		Name:  &createdUser.Name,
		Email: &createdUser.Email,
	}

	return response, nil
}

func (h *Handlers) PatchUsersId(_ context.Context, request User.PatchUsersIdRequestObject) (User.PatchUsersIdResponseObject, error) {
	userID := uint(request.Id)

	existingUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var userToUpdate *userService.Users
	for i, user := range existingUsers {
		if user.ID == userID {
			userToUpdate = &existingUsers[i]
			break
		}
	}

	if userToUpdate == nil {
		return User.PatchUsersId404Response{}, nil
	}
	if request.Body.Name != nil {
		userToUpdate.Name = *request.Body.Name
	}
	if request.Body.Email != nil {
		userToUpdate.Email = *request.Body.Email
	}
	if request.Body.Password != nil {
		userToUpdate.Password = *request.Body.Password
	}

	updatedUser, err := h.Service.UpdateUserByID(userID, *userToUpdate)
	if err != nil {
		return nil, err
	}
	return User.PatchUsersId200JSONResponse{
		Id:       &updatedUser.ID,
		Name:     &updatedUser.Name,
		Email:    &updatedUser.Email,
		Password: &updatedUser.Password,
	}, nil
}

func (h *Handlers) DeleteUsersId(_ context.Context, request User.DeleteUsersIdRequestObject) (User.DeleteUsersIdResponseObject, error) {
	userID := uint(request.Id)

	err := h.Service.DeleteUserByID(userID)
	if err != nil {
		return nil, err
	}
	return User.DeleteUsersId204Response{}, nil
}