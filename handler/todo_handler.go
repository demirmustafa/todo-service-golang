package handler

import (
	"github.com/labstack/echo"
	"todo-service/model"
	error2 "todo-service/model/common/error"
	"todo-service/service"
)

type TodoHandler struct {
	service service.Service
}

type Handler interface {
	HandleGetAll(ctx echo.Context) error
	HandleCreate(ctx echo.Context) error
	HandleUpdate(ctx echo.Context) error
	HandleDelete(ctx echo.Context) error
}

func NewTodoHandler(service service.Service) Handler {
	return TodoHandler{service}
}

func (handler TodoHandler) HandleGetAll(ctx echo.Context) error {
	return ctx.JSON(200, handler.service.FindAll())
}

func (handler TodoHandler) HandleCreate(ctx echo.Context) error {
	var request model.CreateTodoRequest
	_ = ctx.Bind(&request)
	response := handler.service.Create(request)
	error := extractError(response)
	if error != nil {
		return ctx.JSON(error.GetStatus(), error)
	}
	return ctx.JSON(201, response)
}

func (handler TodoHandler) HandleUpdate(ctx echo.Context) error {
	var request model.UpdateTodoRequest
	_ = ctx.Bind(&request)
	id := ctx.Param("id")
	response := handler.service.Update(id, request)
	error := extractError(response)
	if error != nil {
		return ctx.JSON(error.GetStatus(), error)
	}
	return ctx.JSON(202, response)
}

func (handler TodoHandler) HandleDelete(ctx echo.Context) error {
	id := ctx.Param("id")
	handler.service.Delete(id)
	return ctx.NoContent(204)
}

func extractError(i interface{}) error2.Error {
	error, ok := i.(error2.Error)
	if ok {
		return error
	}
	return nil
}
