package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"awesomeProject1/internal/database"
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/taskService"
	"awesomeProject1/internal/web/tasks"
)

func main() {
	database.InitDB() 


	repo := taskService.NewTaskRepository(database.DB)
	service := taskService.NewService(repo)

	handler := handlers.NewHandler(service)
	
	
	e := echo.New()
	
	
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	
	strictHandler := tasks.NewStrictHandler(handler, nil) 
	tasks.RegisterHandlers(e, strictHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}