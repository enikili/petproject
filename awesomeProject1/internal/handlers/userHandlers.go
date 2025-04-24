package handlers

import (
	"context"
	"errors"
	"net/http"

	"awesomeProject1/internal/userService"
	"awesomeProject1/internal/web/users"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
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

    
    response := make([]users.User, len(dbUsers))
    for i, user := range dbUsers {
        userID := int64(user.ID)
        response[i] = users.User{
            Id:       &userID,
            Username: user.Username,
            Email:    types.Email(user.Email),
        }
    }

    return users.GetUsers200JSONResponse(response), nil
}

func (h *UserHandler) PostUsers(ctx context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	if request.Body.Username == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}
	if request.Body.Password == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "password is required")
	}
	if request.Body.Email == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}

	newUser := userService.User{
		Username: request.Body.Username,
		Password: request.Body.Password,
		Email:    string(request.Body.Email),
	}

	createdUser, err := h.Service.CreateUser(newUser)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	userID := int64(createdUser.ID)
	return users.PostUsers201JSONResponse{
		Id:       &userID,
		Username: createdUser.Username,
		Email:    types.Email(createdUser.Email),
	}, nil
}

func (h *UserHandler) PatchUsersById(ctx context.Context, request users.PatchUsersByIdRequestObject) (users.PatchUsersByIdResponseObject, error) {
	userID := uint(request.Id)

	updates := userService.User{}
	if request.Body.Username != "" {
		updates.Username = request.Body.Username
	}
	if request.Body.Password != "" {
		updates.Password = request.Body.Password
	}
	if request.Body.Email != "" {
		updates.Email = string(request.Body.Email)
	}

	updatedUser, err := h.Service.UpdateUser(userID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users.PatchUsersById404JSONResponse{
				Message: "User not found",
			}, nil
		}
		return users.PatchUsersById500JSONResponse{
			Message: "Internal server error",
		}, nil
	}

	userIDResp := int64(updatedUser.ID)
	return users.PatchUsersById200JSONResponse{
		Id:       &userIDResp,
		Username: updatedUser.Username,
		Email:    types.Email(updatedUser.Email),
	}, nil
}
func (h *UserHandler) DeleteUsersById(ctx context.Context, request users.DeleteUsersByIdRequestObject) (users.DeleteUsersByIdResponseObject, error) {
    userID := uint(request.Id)

    if err := h.Service.DeleteUser(userID); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return users.DeleteUsersById404JSONResponse{
                Message: "User not found",
            }, nil
        }
        return users.DeleteUsersById500JSONResponse{
            Message: "Failed to delete user",
        }, nil
    }

    return users.DeleteUsersById204Response{}, nil
}
func (h *UserHandler) GetUsersId(ctx context.Context, request users.GetUsersByIdRequestObject) (users.GetUsersByIdResponseObject, error) {
    userID := uint(request.Id)

    user, err := h.Service.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return users.GetUsersById404JSONResponse{
                Message: "User not found",
            }, nil
        }
        return users.GetUsersById500JSONResponse{
            Message: "Failed to get user",
        }, nil
    }

    userIDResp := int64(user.ID)
    return users.GetUsersById200JSONResponse{
        Id:       &userIDResp,
        Username: user.Username,
        Email:    types.Email(user.Email),
    }, nil
}
func (h *UserHandler) GetUsersById(ctx context.Context, request users.GetUsersByIdRequestObject) (users.GetUsersByIdResponseObject, error) {
    userID := uint(request.Id)

    user, err := h.Service.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return users.GetUsersById404JSONResponse{
                Message: "User not found",
            }, nil
        }
        return users.GetUsersById500JSONResponse{
            Message: "Failed to get user",
        }, nil
    }

    userIDResp := int64(user.ID)
    return users.GetUsersById200JSONResponse{
        Id:       &userIDResp,
        Username: user.Username,
        Email:    types.Email(user.Email),
    }, nil
}