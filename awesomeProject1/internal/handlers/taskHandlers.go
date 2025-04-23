package handlers

import (
	"context"
	"errors"
	"net/http"

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

func (h *TaskHandler) GetTasksId(ctx context.Context, request tasks.GetTasksIdRequestObject) (tasks.GetTasksIdResponseObject, error) {
	userID := uint(request.Id)
	
	dbTasks, err := h.Service.GetTaskByUserId(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get user tasks")
	}

	response := make(tasks.GetTasksId200JSONResponse, 0, len(dbTasks))
	for _, task := range dbTasks {
		taskID := int64(task.ID)
		userID := int64(*task.UserID)
		response = append(response, tasks.Task{
			Id:     &taskID,
			Task:   *task.Task,
			IsDone: task.IsDone,
			UserId: userID,
		})
	}

	return response, nil
}

func (h *TaskHandler) PostTasksId(ctx context.Context, request tasks.PostTasksIdRequestObject) (tasks.PostTasksIdResponseObject, error) {
	userID := uint(request.Id)
	
	if request.Body.Task == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "task description is required")
	}

	newTask := taskService.Task{
		Task:   &request.Body.Task,
		IsDone: request.Body.IsDone,
	}

	createdTask, err := h.Service.CreateTask(userID, newTask)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to create task")
	}

	taskID := int64(createdTask.ID)
	userIDResp := int64(*createdTask.UserID)
	response := tasks.PostTasksId201JSONResponse{
		Id:     &taskID,
		Task:   *createdTask.Task,
		IsDone: createdTask.IsDone,
		UserId: userIDResp,
	}

	return response, nil
}

func (h *TaskHandler) PatchTasksId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
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
			return tasks.PatchTasksId404JSONResponse{
				Message: "Task not found",
			}, nil
		}
		return tasks.PatchTasksId500JSONResponse{
			Message: "Internal server error",
		}, nil
	}

	taskIDResp := int64(updatedTask.ID)
	userIDResp := int64(*updatedTask.UserID)
	return tasks.PatchTasksId200JSONResponse{
		Id:     &taskIDResp,
		Task:   *updatedTask.Task,
		IsDone: updatedTask.IsDone,
		UserId: userIDResp,
	}, nil
}

func (h *TaskHandler) DeleteTasksId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	taskID := uint(request.Id)

	if err := h.Service.DeleteTask(taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.DeleteTasksId404JSONResponse{
				Message: "Task not found",
			}, nil
		}
		return tasks.DeleteTasksId500JSONResponse{
			Message: "Internal server error",
		}, nil
	}

	return tasks.DeleteTasksId204Response{}, nil
}