package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type RequestBody struct {
	Task string `json:"task"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DB.Create(&task)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	DB.Find(&tasks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskByID(taskID uint) (*Task, error) {
	var task Task
	result := DB.First(&task, taskID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	result := DB.Save(task)
	return result.Error
}

func DeleteTask(taskID uint) error {
	result := DB.Delete(&Task{}, taskID)
	return result.Error
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Task ID not found in URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	taskID := uint(id)

	existingTask, err := GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Println("Error getting task:", err)
		return
	}
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if taskValue, ok := updateData["task"].(string); ok {
		existingTask.Task = taskValue
	}
	if isDoneValue, ok := updateData["is_done"].(bool); ok {
		existingTask.IsDone = isDoneValue
	}

	if err := UpdateTask(existingTask); err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		log.Println("Error updating task:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTask)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Task ID not found in URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	taskID := uint(id)

	if err := DeleteTask(taskID); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		log.Println("Error deleting task:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	InitDB()
	DB.AutoMigrate(&Task{})
	router := mux.NewRouter()
	router.HandleFunc("/api/task", CreateTask).Methods("POST")
	router.HandleFunc("/api/task", GetTasks).Methods("GET")
	router.HandleFunc("/api/task/{id:[0-9]+}", UpdateTaskHandler).Methods("PATCH")
	router.HandleFunc("/api/task/{id:[0-9]+}", DeleteTaskHandler).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
