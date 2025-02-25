package handlers

import (
	"awesomeProject1/internal/taskService"
	"awesomeProject1/internal/web/tasks"
	"context"
	"fmt"
	"log"
)

type Handler struct {
	Service *taskService.TaskService
}


func (h *Handler) GetTasks(ctx context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		log.Printf("failed to get all tasks: %v", err)
		return nil, err
	}

	
	response := make(tasks.GetTasks200JSONResponse, 0)

	
	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     &tsk.ID,
			Task:   &tsk.Title,
			IsDone: &tsk.IsDone,
		}
		response = append(response, task)
	}

	
	return response, nil
}


func (h *Handler) PostTasks(ctx context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	
	if request.Body.Task == nil || request.Body.IsDone == nil {
		return nil, fmt.Errorf("task text or is_done field is missing")
	}

	
	taskRequest := request.Body

	
	taskToCreate := taskService.Tasks{
		Title:  *taskRequest.Task,
		IsDone: *taskRequest.IsDone,
	}

	
	createdTask, err := h.Service.CreateTask(taskToCreate)
	if err != nil {
		log.Printf("failed to create task: %v", err)
		return nil, err
	}

	
	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Title,
		IsDone: &createdTask.IsDone,
	}

	
	return response, nil
}


func (h *Handler) PatchTasksId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	
	if request.Body.Task == nil && request.Body.IsDone == nil {
		return nil, fmt.Errorf("at least one field (task or is_done) must be provided")
	}

	
	updates := make(map[string]interface{})
	if request.Body.Task != nil {
		updates["Title"] = *request.Body.Task
	}
	if request.Body.IsDone != nil {
		updates["IsDone"] = *request.Body.IsDone
	}

	
	updatedTask, err := h.Service.PatchTaskByID(uint(request.Id), updates)
	if err != nil {
		log.Printf("failed to update task: %v", err)
		return nil, err
	}

	
	return tasks.PatchTasksId200JSONResponse{
		Id:     &updatedTask.ID,
		Task:   &updatedTask.Title,
		IsDone: &updatedTask.IsDone,
	}, nil
}


func (h *Handler) DeleteTasksId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	
	err := h.Service.DeleteTaskByID(uint(request.Id))
	if err != nil {
		log.Printf("failed to delete task: %v", err)
		return nil, err
	}

	
	return tasks.DeleteTasksId204Response{}, nil
}


func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}