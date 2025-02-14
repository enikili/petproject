package main

import (
	"awesomeProject1/internal/database"
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/taskService"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()
	database.DB.AutoMigrate(&taskService.Task{})

	var repo = taskService.NewTaskRepository(database.DB)
	service := taskService.NewService(repo)

	handler := handlers.NewHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/api/get", handler.GetTasksHandler).Methods("GET")
	router.HandleFunc("/api/post", handler.PostTaskHandler).Methods("POST")
	router.HandleFunc("/api/patch/{id}", handler.UpdateTaskHandler).Methods("PATCH")
	router.HandleFunc("/api/delete/{id}", handler.DeleteTaskHandler).Methods("DELETE")

	log.Println("Starting server on :6050")
	err := http.ListenAndServe(":6050", router)
	if err != nil {
		log.Fatal(err)
	}
}
