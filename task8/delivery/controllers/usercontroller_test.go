package controllers

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

func setupUserRouter(controller *UserController) *gin.Engine {
	router := gin.Default()
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.GET("/users/:email", controller.GetUser)
	router.GET("/users", controller.GetUsers)
	return router
}

func TestRegister(t *testing.T) {
	mockUsecase := new(mocks.UserUsecaseInterface)
	controller := NewUserController(mockUsecase)

	router := setupUserRouter(controller)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("Register", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		reqBody := `{"email": "test@example.com", "password": "password123", "role": "user"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), `"message":"registered successfully"`)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("binding error", func(t *testing.T) {
		reqBody := `invalid json`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"invalid character`)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("Register", mock.AnythingOfType("*domain.User")).Return(errors.New("usecase error")).Once()

		reqBody := `{"email": "test@example.com", "password": "password123", "role": "user"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"usecase error"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUsecase := new(mocks.UserUsecaseInterface)
	controller := NewUserController(mockUsecase)

	router := setupUserRouter(controller)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("Login", mock.AnythingOfType("*domain.User")).Return("validToken", nil).Once()

		reqBody := `{"email": "test@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"token":"validToken"`)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("binding error", func(t *testing.T) {
		reqBody := `invalid json`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"invalid character`)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("Login", mock.AnythingOfType("*domain.User")).Return("", errors.New("invalid credentials")).Once()

		reqBody := `{"email": "test@example.com", "password": "wrongpassword"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"invalid credentials"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mockUsecase := new(mocks.UserUsecaseInterface)
	controller := NewUserController(mockUsecase)

	router := setupUserRouter(controller)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("GetUser", "test@example.com").Return(&domain.User{Email: "test@example.com", Role: "user"}, nil).Once()

		req, _ := http.NewRequest("GET", "/users/test@example.com", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Simulate admin role
		router.GET("/users/:email", func(ctx *gin.Context) {
			ctx.Set("role", "admin")
			controller.GetUser(ctx)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"email":"test@example.com"`)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("forbidden for non-admin role", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/test@example.com", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Simulate non-admin role
		router.GET("/users/:email", func(ctx *gin.Context) {
			ctx.Set("role", "user")
			controller.GetUser(ctx)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), `"message":"admin-only"`)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("GetUser", "test@example.com").Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest("GET", "/users/test@example.com", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Simulate admin role
		router.GET("/users/:email", func(ctx *gin.Context) {
			ctx.Set("role", "admin")
			controller.GetUser(ctx)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"user not found"`)
		mockUsecase.AssertExpectations(t)
	})
}

func TestGetUsers(t *testing.T) {
	mockUsecase := new(mocks.UserUsecaseInterface)
	controller := NewUserController(mockUsecase)

	router := setupUserRouter(controller)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("GetUsers").Return([]domain.User{
			{Email: "user1@example.com", Role: "user"},
			{Email: "user2@example.com", Role: "user"},
		}, nil).Once()

		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Simulate admin role
		router.GET("/users", func(ctx *gin.Context) {
			ctx.Set("role", "admin")
			controller.GetUsers(ctx)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"email":"user1@example.com"`)
		assert.Contains(t, rr.Body.String(), `"email":"user2@example.com"`)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("forbidden for non-admin role", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Simulate non-admin role
		router.GET("/users", func(ctx *gin.Context) {
			ctx.Set("role", "user")
			controller.GetUsers(ctx)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), `"message":"admin-only"`)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("GetUsers").Return(nil, errors.New("could not fetch users")).Once()

		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// Simulate admin role
		router.GET("/users", func(ctx *gin.Context) {
			ctx.Set("role", "admin")
			controller.GetUsers(ctx)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), `"error":"could not fetch users"`)
		mockUsecase.AssertExpectations(t)
	})
}
