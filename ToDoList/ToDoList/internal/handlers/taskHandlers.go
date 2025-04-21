package handlers

import (
	"awesomeProject1/internal/taskService"
	"awesomeProject1/internal/web/tasks"
	"context"
	"net/http"
	"github.com/labstack/echo/v4"
)
type Handler struct {
	Service *taskService.TaskService
}

func NewTaskHandler(service *taskService.TaskService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var response []tasks.Task

	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     &tsk.ID,
			Task:   *tsk.Task,
			IsDone: tsk.IsDone,
		}
		response = append(response, task)
	}

	return tasks.GetTasks200JSONResponse(response), nil
}

func (h *Handler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	taskRequest := request.Body

	taskToCreate := taskService.Task{
		Task:   &taskRequest.Task,
		IsDone: taskRequest.IsDone,
		UserID: taskRequest.Id,
	}
	createdTask, err := h.Service.CreateTask(taskToCreate)
	if err != nil {
		return nil, err
	}

	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   *createdTask.Task,
		IsDone: createdTask.IsDone,
		UserId: int(*createdTask.UserID),
	}

	return response, nil
}

func (h *Handler) PatchTasksId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
    
    updateData := taskService.Task{
        UserID: request.Body.Id,
    }

    
    if request.Body.Task != "" { 
        updateData.Task = &request.Body.Task
    }
    if request.Body.IsDone != nil { 
        updateData.IsDone = *&request.Body.IsDone
    }
    if request.Body.UserId != 0 { 
        updateData.UserID = request.Body.Id
    }

    updatedTask, err := h.Service.UpdateTaskByID(uint(request.Id), updateData)
    if err != nil {
        return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to update task")
    }

    taskID := updatedTask.ID
    return tasks.PatchTasksId200JSONResponse{
        Id:     &taskID,
        Task:   *updatedTask.Task,
        IsDone: updatedTask.IsDone,
        UserId: int(*updatedTask.UserID),
    }, nil
}

func (h *Handler) DeleteTasksId(_ context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	taskID := uint(request.Id)

	err := h.Service.DeleteTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	return tasks.DeleteTasksId204Response{}, nil
}