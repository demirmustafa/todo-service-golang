package service

import (
	"todo-service/model"
	err "todo-service/model/common/error"
	"todo-service/repository"
)

type Service interface {
	FindAll() interface{}
	Find(id string) interface{}
	Create(i interface{}) interface{}
	Update(id string, i interface{}) interface{}
	Delete(id string)
}

type TodoService struct {
	repository repository.Repository
}

func NewTodoService(repository repository.Repository) Service {
	return TodoService{repository}
}

func (s TodoService) FindAll() interface{} {
	return s.repository.FindAll()
}

func (s TodoService) Find(id string) interface{} {
	found := s.repository.Find(id)
	if found == nil {
		return err.NotFoundError{Message: "No todo found with provided id!"}
	}
	return found
}

func (s TodoService) Create(i interface{}) interface{} {
	request, ok := i.(model.CreateTodoRequest)
	if !ok {
		return err.BadRequestError{Message: "Invalid request!"}
	}

	if request.Title == "" {
		return err.BadRequestError{Message: "Title cannot be empty!"}
	}

	todo := model.Todo{
		Title:     request.Title,
		Completed: false,
	}
	return s.repository.Create(todo)
}

func (s TodoService) Update(id string, i interface{}) interface{} {
	request, ok := i.(model.UpdateTodoRequest)
	if !ok {
		return err.BadRequestError{Message: "Invalid request!"}
	}

	if request.Title == "" {
		return err.BadRequestError{Message: "Title cannot be empty"}
	}

	found := s.repository.Find(id)
	if found == nil {
		return err.NotFoundError{Message: "No todo found with provided id!"}
	}
	todo := found.(model.Todo)
	todo.Title = request.Title
	todo.Completed = request.Completed
	return s.repository.Update(todo)
}

func (s TodoService) Delete(id string) {
	s.repository.Delete(id)
}
