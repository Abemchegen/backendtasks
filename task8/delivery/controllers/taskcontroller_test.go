package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"task8/domain"
	"task8/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupRouter initializes the Gin router with the necessary routes for testing.
func setupRouter(taskUsecase domain.TaskUsecaseInterface) *gin.Engine {
	router := gin.Default()
	router.POST("/tasks", func(ctx *gin.Context) {
		// Simulate user role and user ID for testing
		ctx.Set("role", "user")
		ctx.Set("user_id", "userID")

		var newTask domain.Task
		if err := ctx.ShouldBindJSON(&newTask); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := ctx.GetString("user_id")
		err := taskUsecase.CreateTask(&newTask, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "task created"})
	})

	router.GET("/tasks/:id", func(ctx *gin.Context) {
		// Simulate user role for testing
		ctx.Set("role", "user")

		taskID := ctx.Param("id")
		task, err := taskUsecase.GetTask(taskID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, task)
	})

	router.GET("/tasks", func(ctx *gin.Context) {
		// Simulate user role and user ID for testing
		ctx.Set("role", "user")
		ctx.Set("user_id", "userID")

		userID := ctx.GetString("user_id")
		tasks, err := taskUsecase.GetTasks(userID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, tasks)
	})

	router.PUT("/tasks/:id", func(ctx *gin.Context) {
		// Simulate user role and user ID for testing
		ctx.Set("role", "user")
		ctx.Set("user_id", "userID")

		var updatedTask domain.Task
		if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		taskID := ctx.Param("id")
		// userID := ctx.GetString("user_id")
		err := taskUsecase.UpdateTask(taskID, &updatedTask)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "task updated"})
	})

	router.DELETE("/tasks/:id", func(ctx *gin.Context) {
		// Simulate user role and user ID for testing
		ctx.Set("role", "user")
		ctx.Set("user_id", "userID")

		taskID := ctx.Param("id")
		err := taskUsecase.RemoveTask(taskID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "task deleted"})
	})

	return router
}

func TestCreateTask(t *testing.T) {
	mockUsecase := new(mocks.TaskUsecaseInterface)

	router := setupRouter(mockUsecase)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("CreateTask", mock.AnythingOfType("*domain.Task"), "userID").Return(nil).Once()

		reqBody := `{"title": "New Task", "description": "Task Description"}`
		req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("binding error", func(t *testing.T) {
		reqBody := `invalid json`
		req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"invalid character`)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("CreateTask", mock.AnythingOfType("*domain.Task"), "userID").Return(errors.New("usecase error")).Once()

		reqBody := `{"title": "New Task", "description": "Task Description"}`
		req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"usecase error"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestGetTask(t *testing.T) {
	mockUsecase := new(mocks.TaskUsecaseInterface)

	router := setupRouter(mockUsecase)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("GetTask", "taskID").Return(&domain.Task{Title: "Task Title", Description: "Task Description"}, nil).Once()

		req, _ := http.NewRequest("GET", "/tasks/taskID", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"title":"Task Title"`)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("GetTask", "taskID").Return(nil, errors.New("task not found")).Once()

		req, _ := http.NewRequest("GET", "/tasks/taskID", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"task not found"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestGetTasks(t *testing.T) {
	mockUsecase := new(mocks.TaskUsecaseInterface)

	router := setupRouter(mockUsecase)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("GetTasks", "userID").Return(&[]domain.Task{
			{Title: "Task 1", Description: "Description 1"},
			{Title: "Task 2", Description: "Description 2"},
		}, nil).Once()

		req, _ := http.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"title":"Task 1"`)
		assert.Contains(t, rr.Body.String(), `"title":"Task 2"`)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("GetTasks", "userID").Return(nil, errors.New("could not fetch tasks")).Once()

		req, _ := http.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"could not fetch tasks"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	mockUsecase := new(mocks.TaskUsecaseInterface)

	router := setupRouter(mockUsecase)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("UpdateTask", "taskID", mock.AnythingOfType("*domain.Task")).Return(nil).Once()

		reqBody := `{"title": "Updated Task", "description": "Updated Description"}`
		req, _ := http.NewRequest("PUT", "/tasks/taskID", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("binding error", func(t *testing.T) {
		reqBody := `invalid json`
		req, _ := http.NewRequest("PUT", "/tasks/taskID", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"invalid character`)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("UpdateTask", "taskID", mock.AnythingOfType("*domain.Task")).Return(errors.New("usecase error")).Once()

		reqBody := `{"title": "Updated Task", "description": "Updated Description"}`
		req, _ := http.NewRequest("PUT", "/tasks/taskID", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"usecase error"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestRemoveTask(t *testing.T) {
	mockUsecase := new(mocks.TaskUsecaseInterface)

	router := setupRouter(mockUsecase)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("RemoveTask", "taskID").Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/tasks/taskID", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("RemoveTask", "taskID").Return(errors.New("usecase error")).Once()

		req, _ := http.NewRequest("DELETE", "/tasks/taskID", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"usecase error"`)
		mockUsecase.AssertExpectations(t)
	})
}
