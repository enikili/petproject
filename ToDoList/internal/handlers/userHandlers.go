package handlers

import (
    "awesomeProject1/internal/userService"
    "github.com/labstack/echo/v4"
    "net/http"
    "strconv"
)

type UserHandler struct {
    Service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
    return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
    var input userService.UserCreate
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    user, err := h.Service.CreateUser(&input)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
    users, err := h.Service.GetAllUsers()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
    idStr := c.Param("id") // Извлекаем ID из URL
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
    }

    user, err := h.Service.GetUserByID(uint(id))
    if err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
    }

    return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
    }

    var input userService.UserUpdate
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    user, err := h.Service.UpdateUser(uint(id), &input)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
    }

    if err := h.Service.DeleteUser(uint(id)); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.NoContent(http.StatusNoContent)
}