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

// GetTasks возвращает список всех задач.
func (h *Handler) GetTasks(ctx context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	// Получение всех задач из сервиса
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		log.Printf("failed to get all tasks: %v", err)
		return nil, err
	}

	// Создаем пустой срез для ответа
	response := make(tasks.GetTasks200JSONResponse, 0)

	// Заполняем слайс response всеми задачами из БД
	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     &tsk.ID,
			Task:   &tsk.Title,
			IsDone: &tsk.IsDone,
		}
		response = append(response, task)
	}

	// Возвращаем ответ
	return response, nil
}

// PostTasks создает новую задачу.
func (h *Handler) PostTasks(ctx context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	// Проверяем, что поля Task и IsDone не nil
	if request.Body.Task == nil || request.Body.IsDone == nil {
		return nil, fmt.Errorf("task text or is_done field is missing")
	}

	// Распаковываем тело запроса
	taskRequest := request.Body

	// Создаем задачу для добавления в сервис
	taskToCreate := taskService.Tasks{
		Title:  *taskRequest.Task,
		IsDone: *taskRequest.IsDone,
	}

	// Создаем задачу в сервисе
	createdTask, err := h.Service.CreateTask(taskToCreate)
	if err != nil {
		log.Printf("failed to create task: %v", err)
		return nil, err
	}

	// Создаем структуру ответа
	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Title,
		IsDone: &createdTask.IsDone,
	}

	// Возвращаем ответ
	return response, nil
}

// PatchTasksId обновляет задачу по её ID.
func (h *Handler) PatchTasksId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	// Проверяем, что хотя бы одно поле для обновления передано
	if request.Body.Task == nil && request.Body.IsDone == nil {
		return nil, fmt.Errorf("at least one field (task or is_done) must be provided")
	}

	// Создаем карту для обновления
	updates := make(map[string]interface{})
	if request.Body.Task != nil {
		updates["Title"] = *request.Body.Task
	}
	if request.Body.IsDone != nil {
		updates["IsDone"] = *request.Body.IsDone
	}

	// Обновляем задачу в сервисе
	updatedTask, err := h.Service.PatchTaskByID(uint(request.Id), updates)
	if err != nil {
		log.Printf("failed to update task: %v", err)
		return nil, err
	}

	// Возвращаем обновлённую задачу
	return tasks.PatchTasksId200JSONResponse{
		Id:     &updatedTask.ID,
		Task:   &updatedTask.Title,
		IsDone: &updatedTask.IsDone,
	}, nil
}

// DeleteTasksId удаляет задачу по её ID.
func (h *Handler) DeleteTasksId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	// Удаляем задачу в сервисе
	err := h.Service.DeleteTaskByID(uint(request.Id))
	if err != nil {
		log.Printf("failed to delete task: %v", err)
		return nil, err
	}

	// Возвращаем статус 204 No Content
	return tasks.DeleteTasksId204Response{}, nil
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}