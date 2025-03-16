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

    
    Taskrepo := taskService.NewTaskRepository(database.DB)
    Taskservice := taskService.NewService(Taskrepo)
    taskHandler := handlers.NewHandler(Taskservice)

    
    UserRepo := userService.NewUserRepository(database.DB)
    UserService := userService.NewUserService(UserRepo)
    UserHandler := handlers.NewUserHandler(UserService)

    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    
    strictHandler := tasks.NewStrictHandler(taskHandler, nil)
    tasks.RegisterHandlers(e, strictHandler)

    // Регистрация маршрутов для пользователей
    e.POST("/users", UserHandler.CreateUser)          
    e.GET("/users", UserHandler.GetAllUsers)          
    e.GET("/users/:id", UserHandler.GetUserByID)      
    e.PUT("/users/:id", UserHandler.UpdateUser)       
    e.DELETE("/users/:id", UserHandler.DeleteUser)    

    if err := e.Start(":8080"); err != nil {
        log.Fatalf("failed to start with err: %v", err)
    }
}