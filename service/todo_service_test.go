package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"todo-service/model"
	err "todo-service/model/common/error"
	"todo-service/repository"
)

func Test_ShouldFindAllTodos(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var todoList []model.Todo
	todoList = append(todoList, model.Todo{
		ID:        "12345",
		Title:     "test todo",
		Completed: false,
	})

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().FindAll().Return(todoList).Times(1)
	service := TodoService{repository}

	all := service.FindAll().([]model.Todo)
	assert.NotNil(t, all)
	assert.NotEmpty(t, all)
	assert.Equal(t, 1, len(all))
	assert.Equal(t, "12345", all[0].ID)
	assert.Equal(t, "test todo", all[0].Title)
	assert.False(t, all[0].Completed)
}

func Test_ShouldFindTodoById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	expected := model.Todo{
		ID:        "12345",
		Title:     "test todo",
		Completed: true,
	}

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().Find(gomock.Eq("12345")).Return(expected).Times(1)

	service := NewTodoService(repository)
	todo := service.Find("12345").(model.Todo)

	assert.NotNil(t, todo)
	assert.Equal(t, "12345", todo.ID)
	assert.Equal(t, "test todo", todo.Title)
	assert.True(t, todo.Completed)
}

func Test_ShouldReturnNotFoundErrorWhenNoTodoFoundById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().Find(gomock.Eq("12345")).Return(nil).Times(1)

	service := NewTodoService(repository)
	found := service.Find("12345")

	assert.NotNil(t, found)
	assert.IsType(t, err.NotFoundError{}, found)
	err := found.(err.Error)
	assert.Equal(t, 404, err.GetStatus())
	assert.Equal(t, "No todo found with provided id!", err.GetMessage())
}

func Test_ShouldCreateTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().Create(model.Todo{Title: "test todo", Completed: false}).Return(model.Todo{
		ID:        "12345",
		Title:     "test todo",
		Completed: false,
	}).Times(1)

	service := NewTodoService(repository)
	created := service.Create(model.CreateTodoRequest{Title: "test todo"})

	assert.NotNil(t, created)
	todo := created.(model.Todo)
	assert.Equal(t, "12345", todo.ID)
	assert.Equal(t, "test todo", todo.Title)
	assert.False(t, todo.Completed)
}

func Test_ShouldReturnBadRequestErrorWhenRequestIsEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := TodoService{nil}
	response := service.Create(nil)

	error, ok := response.(err.BadRequestError)
	if !ok {
		assert.Fail(t, "Expected BadRequestError response")
	}

	assert.NotNil(t, error)
	assert.Equal(t, 400, error.GetStatus())
	assert.Equal(t, "Invalid request!", error.GetMessage())
}

func Test_ShouldReturnBadRequestErrorWhenTitleIsEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := TodoService{nil}
	response := service.Create(model.CreateTodoRequest{})

	error, ok := response.(err.BadRequestError)
	if !ok {
		assert.Fail(t, "Expected BadRequestError response")
	}

	assert.NotNil(t, error)
	assert.Equal(t, 400, error.GetStatus())
	assert.Equal(t, "Title cannot be empty!", error.GetMessage())
}

func Test_ShouldUpdateTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	beforeUpdate := model.Todo{
		ID:        "12345",
		Title:     "test todo before update",
		Completed: false,
	}

	afterUpdate := model.Todo{
		ID:        "12345",
		Title:     "test todo after update",
		Completed: true,
	}

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().Find(gomock.Eq("12345")).Return(beforeUpdate).Times(1)
	repository.EXPECT().Update(gomock.Eq(afterUpdate)).Return(afterUpdate).Times(1)

	service := TodoService{repository}
	updated := service.Update("12345", model.UpdateTodoRequest{
		Title:     "test todo after update",
		Completed: true,
	})

	assert.NotNil(t, updated)

	todo, ok := updated.(model.Todo)
	assert.True(t, ok)
	assert.Equal(t, "12345", todo.ID)
	assert.Equal(t, "test todo after update", todo.Title)
	assert.True(t, todo.Completed)
}

func Test_ShouldReturnBadRequestErrorWhenUpdateRequestIsEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := TodoService{nil}
	response := service.Update("12345", nil)

	error, ok := response.(err.BadRequestError)
	assert.True(t, ok)
	assert.Equal(t, 400, error.GetStatus())
	assert.Equal(t, "Invalid request!", error.GetMessage())
}

func Test_ShouldReturnBadRequestErrorWhenUpdateRequestTitleIsEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := TodoService{nil}
	response := service.Update("12345", model.UpdateTodoRequest{
		Title:     "",
		Completed: false,
	})

	error, ok := response.(err.BadRequestError)
	assert.True(t, ok)
	assert.Equal(t, 400, error.GetStatus())
	assert.Equal(t, "Title cannot be empty", error.GetMessage())
}

func Test_ShouldReturnNotFoundErrorWhenNoTodoFoundToUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().Find(gomock.Eq("12345")).Return(nil).Times(1)

	service := TodoService{repository}
	response := service.Update("12345", model.UpdateTodoRequest{
		Title:     "test todo",
		Completed: true,
	})

	error, ok := response.(err.NotFoundError)
	assert.True(t, ok)
	assert.Equal(t, 404, error.GetStatus())
	assert.Equal(t, "No todo found with provided id!", error.GetMessage())
}

func Test_ShouldDeleteTodo(t *testing.T) {
	controller  := gomock.NewController(t)
	defer controller.Finish()

	repository := repository.NewMockRepository(controller)
	repository.EXPECT().Delete(gomock.Eq("12345")).Times(1)

	service := TodoService{repository}
	service.Delete("12345")
}