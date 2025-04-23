package main

import (
	"awesomeProject1/internal/database"
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/taskService"
	"awesomeProject1/internal/userService"
	"awesomeProject1/internal/web/tasks"
	"awesomeProject1/internal/web/users"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.InitDB()
	

	
	tasksRepo := taskService.NewTaskRepository(database.DB)
	usersRepo := userService.NewUserRepository(database.DB)
	tasksService := taskService.NewService(tasksRepo)
	usersService := userService.NewUserService(usersRepo)

	
	tasksHandler := handlers.NewTaskHandler(tasksService) 
	usersHandler := handlers.NewUserHandler(usersService) 


	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	
	strictTasksHandler := tasks.NewStrictHandler(tasksHandler, nil)
	strictUsersHandler := users.NewStrictHandler(usersHandler, nil)

	
	tasks.RegisterHandlers(e, strictTasksHandler)
	users.RegisterHandlers(e, strictUsersHandler)
	
	for _, route := range e.Routes() {
    fmt.Printf("Registered route: %6s %s\n", route.Method, route.Path)
}
	
	log.Println("Server starting on :8084...")
	if err := e.Start(":8084"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}