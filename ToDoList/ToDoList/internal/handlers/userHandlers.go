package handlers

import (
	"context"
	"errors"
	"net/http"
	"awesomeProject1/internal/web/users"
	"awesomeProject1/internal/userService"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	Service *userService.Service
}

func NewUserHandler(service *userService.Service) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) GetUsers(ctx context.Context, request api.GetUsersRequestObject) (api.GetUsersResponseObject, error) {
	allUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get users: "+err.Error())
	}

	response := make(api.GetUsers200JSONResponse, 0, len(allUsers))
	for _, user := range allUsers {
		userID := user.ID
		response = append(response, api.User{
			Id:    &userID,
			Email: &user.Email,
			Name:  &user.Name,
		})
	}

	return response, nil
}

func (h *UserHandler) GetUsersIdTasks(ctx context.Context, request api.GetUsersIdTasksRequestObject) (api.GetUsersIdTasksResponseObject, error) {
	userID := uint(request.Id)

	dbTasks, err := h.Service.GetTasksForUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get user tasks: "+err.Error())
	}

	response := make(api.GetUsersIdTasks200JSONResponse, 0, len(dbTasks))
	for _, task := range dbTasks {
		taskID := task.ID
		userID := task.UserID
		response = append(response, api.Task{
			Id:     &taskID,
			Task:   task.Task,
			IsDone: task.IsDone,
			UserId: &userID,
		})
	}

	return response, nil
}

func (h *UserHandler) PostUsers(ctx context.Context, request api.PostUsersRequestObject) (api.PostUsersResponseObject, error) {
	if request.Body.Email == nil || *request.Body.Email == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}
	if request.Body.Password == nil || *request.Body.Password == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "password is required")
	}

	newUser := userService.User{
		Email:    *request.Body.Email,
		Name:     request.Body.Name, // Name может быть nil
		Password: *request.Body.Password,
	}

	createdUser, err := h.Service.CreateUser(newUser)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "failed to create user: "+err.Error())
	}

	response := api.PostUsers201JSONResponse{
		Id:    &createdUser.ID,
		Name:  &createdUser.Name,
		Email: &createdUser.Email,
	}

	return response, nil
}

func (h *UserHandler) PatchUsersId(ctx context.Context, request api.PatchUsersIdRequestObject) (api.PatchUsersIdResponseObject, error) {
	userID := uint(request.Id)

	updates := userService.User{}
	if request.Body.Name != nil {
		updates.Name = *request.Body.Name
	}
	if request.Body.Email != nil {
		if *request.Body.Email == "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "email cannot be empty")
		}
		updates.Email = *request.Body.Email
	}
	if request.Body.Password != nil {
		if *request.Body.Password == "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "password cannot be empty")
		}
		updates.Password = *request.Body.Password
	}

	updatedUser, err := h.Service.UpdateUserByID(userID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.PatchUsersId404Response{}, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to update user: "+err.Error())
	}

	return api.PatchUsersId200JSONResponse{
		Id:    &updatedUser.ID,
		Name:  &updatedUser.Name,
		Email: &updatedUser.Email,
	}, nil
}

func (h *UserHandler) DeleteUsersId(ctx context.Context, request api.DeleteUsersIdRequestObject) (api.DeleteUsersIdResponseObject, error) {
	userID := uint(request.Id)

	if err := h.Service.DeleteUserByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.DeleteUsersId404Response{}, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user: "+err.Error())
	}

	return api.DeleteUsersId204Response{}, nil
}