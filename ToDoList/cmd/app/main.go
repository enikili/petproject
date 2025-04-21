package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "log"
    "awesomeProject1/internal/database"
    "awesomeProject1/internal/handlers"
    "awesomeProject1/internal/taskService"
    "awesomeProject1/internal/userService"
    "awesomeProject1/internal/web/tasks"
)

func main() {

	database.InitDB()

	repoTask := taskService.NewTaskRepository(database.DB)
	serviceTask := taskService.NewTaskService(repoTask)

	repoUser := userService.NewUserRepository(database.DB, repoTask)
	serviceUser := userService.NewUserService(repoUser)

	taskHTTPHandler := handlers.NewTaskHandler(serviceTask)
	userHTTPHandler := handlers.NewUserHandler(serviceUser)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	tasksHandler := tasks.NewStrictHandler(taskHTTPHandler, nil)
	tasks.RegisterHandlers(e, tasksHandler)

	usersHandler := User.NewStrictHandler(userHTTPHandler, nil)
	User.RegisterHandlers(e, usersHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatal("failed to start with error", err)
	}
}