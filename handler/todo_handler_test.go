package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-service/model"
	error2 "todo-service/model/common/error"
	"todo-service/service"
)

func Test_ShouldReturnEmptyTodoList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	todos := make([]model.Todo, 0)
	service := service.NewMockService(controller)
	service.EXPECT().FindAll().Return(todos).Times(1)

	handler := TodoHandler{service}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	expected := `[]
`
	if assert.NoError(t, handler.HandleGetAll(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func Test_ShouldFindAllTodos(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var todos []model.Todo
	todos = append(todos, model.Todo{
		ID:        "12345",
		Title:     "test todo",
		Completed: false,
	})

	service := service.NewMockService(controller)
	service.EXPECT().FindAll().Return(todos).Times(1)
	handler := TodoHandler{service}

	expected := `[{"id":"12345","title":"test todo","completed":false}]
`

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.HandleGetAll(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func Test_ShouldCreateTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	request := model.CreateTodoRequest{Title: "test todo"}
	requestAsString := `{"title":"test todo"}`
	expected := model.Todo{
		ID:        "12345",
		Title:     "test todo",
		Completed: false,
	}
	expectedAsString := `{"id":"12345","title":"test todo","completed":false}
`

	service := service.NewMockService(controller)
	service.EXPECT().Create(gomock.Eq(request)).Return(expected).Times(1)

	handler := TodoHandler{service}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(requestAsString))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.HandleCreate(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedAsString, rec.Body.String())
	}
}

func Test_ShouldReturnBadRequestErrorWhenCreatingTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	request := model.CreateTodoRequest{}
	requestAsString := `{}`
	expected := error2.BadRequestError{Message: "test bad request error message"}
	expectedAsString := `{"message":"test bad request error message"}
`
	service := service.NewMockService(controller)
	service.EXPECT().Create(gomock.Eq(request)).Return(expected).Times(1)

	handler := TodoHandler{service}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(requestAsString))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.HandleCreate(ctx)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedAsString, rec.Body.String())
	}
}

func Test_ShouldUpdateTodoStatusAsCompleted(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	request := model.UpdateTodoRequest{
		Title:     "test todo",
		Completed: true,
	}

	requestAsString := `{"title":"test todo","completed":true}
`
	expected := model.Todo{
		ID:        "12345",
		Title:     "test todo",
		Completed: true,
	}
	expectedAsString := `{"id":"12345","title":"test todo","completed":true}
`
	service := service.NewMockService(controller)
	service.EXPECT().Update(gomock.Eq("12345"), gomock.Eq(request)).Return(expected).Times(1)
	handler := TodoHandler{service}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/todos/12345", strings.NewReader(requestAsString))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("12345")

	if assert.NoError(t, handler.HandleUpdate(ctx)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, expectedAsString, rec.Body.String())
	}
}

func Test_ShouldReturnBadRequestErrorWhenUpdatingTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	request := model.UpdateTodoRequest{}
	requestAsString := `{}`
	expected := error2.BadRequestError{Message: "test bad request error message"}
	expectedAsString := `{"message":"test bad request error message"}
`
	service := service.NewMockService(controller)
	service.EXPECT().Update(gomock.Eq("12345"), gomock.Eq(request)).Return(expected).Times(1)

	handler := TodoHandler{service}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/todos/12345", strings.NewReader(requestAsString))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("12345")

	if assert.NoError(t, handler.HandleUpdate(ctx)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedAsString, rec.Body.String())
	}
}

func Test_ShouldDeleteTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	service := service.NewMockService(controller)
	service.EXPECT().Delete(gomock.Eq("12345"))
	handler := TodoHandler{service}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/12345", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("12345")

	if assert.NoError(t, handler.HandleDelete(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
