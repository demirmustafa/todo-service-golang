package main

import (
	"github.com/labstack/echo"
	"os"
	"todo-service/configuration"
	"todo-service/handler"
	"todo-service/repository"
	"todo-service/service"
)

func main() {
	e := echo.New()

	config := new(configuration.MongoConfiguration).Init(getDBUri(), getDBName())

	todoRepository := repository.NewTodoRepository(config.Database().Collection("todos"))
	todoService := service.NewTodoService(todoRepository)
	todoHandler := handler.NewTodoHandler(todoService)
	e.GET("/todos", todoHandler.HandleGetAll)
	e.POST("/todos", todoHandler.HandleCreate)
	e.PUT("/todos/:id", todoHandler.HandleUpdate)
	e.DELETE("/todos/:id", todoHandler.HandleDelete)

	e.Logger.Fatal(e.Start(":8080"))
}

func getDBUri() string {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return "mongodb://localhost:27017"
	}
	return uri
}

func getDBName() string {
	return "todo-service"
}

