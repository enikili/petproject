package handlers

import (
	"context"
	"errors"
	"net/http"
	"awesomeProject1/internal/userService"
	"awesomeProject1/internal/web/users"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	Service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) GetUsers(ctx context.Context, request users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	dbUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get users")
	}

	response := make(users.GetUsers200JSONResponse, 0, len(dbUsers))
	for _, user := range dbUsers {
		userID := int64(user.ID)
		response = append(response, users.User{
			Id:       &userID,
			Username: user.Username,
		})
	}

	return response, nil
}

func (h *UserHandler) PostUsers(ctx context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	if request.Body.Username == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}
	if request.Body.Password == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "password is required")
	}
	if request.Body.Email == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}

	newUser := userService.User{
		Username: *request.Body.Username,
		Password: *request.Body.Password,
		Email:    string(*request.Body.Email),
	}

	createdUser, err := h.Service.CreateUser(newUser)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	userID := int64(createdUser.ID)
	return users.PostUsers201JSONResponse{
		Id:       &userID,
		Username: createdUser.Username,
	}, nil
}

func (h *UserHandler) PatchUsersId(ctx context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
    userID := uint(request.Id)
    
    updates := userService.User{}
    if request.Body.Username != nil {
        updates.Username = *request.Body.Username
    }
    if request.Body.Password != nil {
        updates.Password = *request.Body.Password
    }
    if request.Body.Email != nil {
        
        updates.Email = string(*request.Body.Email)
    }

    updatedUser, err := h.Service.UpdateUserByID(userID, updates)
    if err != nil {
    return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

    userIDResp := int64(updatedUser.ID)
    return users.PatchUsersId200JSONResponse{
        Id:       &userIDResp,
        Username: updatedUser.Username,
        Email:    updatedUser.Email, 
    }, nil
}

func (h *UserHandler) DeleteUsersId(ctx context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	userID := uint(request.Id)

	if err := h.Service.DeleteUserByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users.DeleteUsersId404JSONResponse{
				Message: "User not found",
			}, nil
		}
		return users.DeleteUsersId500JSONResponse{
			Message: "Internal server error",
		}, nil
	}

	return users.DeleteUsersId204Response{}, nil
}