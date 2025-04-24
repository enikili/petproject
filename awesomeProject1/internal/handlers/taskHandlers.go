package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"awesomeProject1/internal/taskService"
	"awesomeProject1/internal/web/tasks"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TaskHandler struct {
	Service *taskService.TaskService
}

func NewTaskHandler(service *taskService.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

func (h *TaskHandler) GetTasks(ctx context.Context, request tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get tasks")
	}

	response := make([]tasks.Task, 0, len(allTasks))
	for _, t := range allTasks {
		taskID := int64(t.ID)
		userID := int64(*t.UserID)
		response = append(response, tasks.Task{
			Id:     &taskID,
			Task:   *t.Task,
			IsDone: t.IsDone,
			UserId: userID,
		})
	}

	return tasks.GetTasks200JSONResponse(response), nil
}

func (h *TaskHandler) GetTasksById(ctx context.Context, request tasks.GetTasksByIdRequestObject) (tasks.GetTasksByIdResponseObject, error) {
	taskID := uint(request.Id)

	task, err := h.Service.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.GetTasksById404JSONResponse{
				Message: "Task not found",
			}, nil
		}
		return tasks.GetTasksById500JSONResponse{
			Message: "Failed to get task",
		}, nil
	}

	taskIDResp := int64(task.ID)
	userIDResp := int64(*task.UserID)
	return tasks.GetTasksById200JSONResponse{
		Id:     &taskIDResp,
		Task:   *task.Task,
		IsDone: task.IsDone,
		UserId: userIDResp,
	}, nil
}

func (h *TaskHandler) GetUsersIdTasks(ctx context.Context, request tasks.GetUsersIdTasksRequestObject) (tasks.GetUsersIdTasksResponseObject, error) {
	userID := uint(request.Id)
	
	dbTasks, err := h.Service.GetTasksId(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.GetUsersIdTasks404JSONResponse{
				Message: "User not found",
			}, nil
		}
		return tasks.GetUsersIdTasks500JSONResponse{
			Message: "Failed to get user tasks",
		}, nil
	}

	response := make([]tasks.Task, len(dbTasks))
	for i, task := range dbTasks {
		taskID := int64(task.ID)
		userID := int64(*task.UserID)
		response[i] = tasks.Task{
			Id:     &taskID,
			Task:   *task.Task,
			IsDone: task.IsDone,
			UserId: userID,
		}
	}

	return tasks.GetUsersIdTasks200JSONResponse(response), nil
}

func (h *TaskHandler) PostTasks(ctx context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	if request.Body.Task == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "task description is required")
	}
	if request.Body.UserId == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	userID := uint(request.Body.UserId)
	newTask := taskService.Task{
		Task:   &request.Body.Task,
		IsDone: request.Body.IsDone,
		UserID: &userID,
	}

	createdTask, err := h.Service.CreateTask(userID, newTask)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return nil, echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to create task")
	}

	taskID := int64(createdTask.ID)
	userIDResp := int64(*createdTask.UserID)
	return tasks.PostTasks201JSONResponse{
		Id:     &taskID,
		Task:   *createdTask.Task,
		IsDone: createdTask.IsDone,
		UserId: userIDResp,
	}, nil
}

func (h *TaskHandler) PatchTasksById(ctx context.Context, request tasks.PatchTasksByIdRequestObject) (tasks.PatchTasksByIdResponseObject, error) {
	taskID := uint(request.Id)

	updates := taskService.Task{}
	if request.Body.Task != "" {
		updates.Task = &request.Body.Task
	}
	if request.Body.IsDone != nil {
		updates.IsDone = request.Body.IsDone
	}

	updatedTask, err := h.Service.UpdateTaskByID(taskID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.PatchTasksById404JSONResponse{
				Message: "Task not found",
			}, nil
		}
		return tasks.PatchTasksById500JSONResponse{
			Message: "Internal server error",
		}, nil
	}

	taskIDResp := int64(updatedTask.ID)
	userIDResp := int64(*updatedTask.UserID)
	return tasks.PatchTasksById200JSONResponse{
		Id:     &taskIDResp,
		Task:   *updatedTask.Task,
		IsDone: updatedTask.IsDone,
		UserId: userIDResp,
	}, nil
}

func (h *TaskHandler) DeleteTasksById(ctx context.Context, request tasks.DeleteTasksByIdRequestObject) (tasks.DeleteTasksByIdResponseObject, error) {
	taskID := uint(request.Id)

	if err := h.Service.DeleteTasksById(taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.DeleteTasksById404JSONResponse{
				Message: "Task not found",
			}, nil
		}
		return tasks.DeleteTasksById500JSONResponse{
			Message: "Internal server error",
		}, nil
	}

	return tasks.DeleteTasksById204Response{}, nil
}