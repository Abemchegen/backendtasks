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

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestUserController_Register(t *testing.T) {
	router := setupRouter()
	mockUserUsecase := new(mocks.UserUsecaseInterface)
	userController := NewUserController(mockUserUsecase)

	router.POST("/register", userController.Register)

	t.Run("successful registration", func(t *testing.T) {
		mockUserUsecase.On("Register", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		w := httptest.NewRecorder()
		reqBody := `{"email":"test@example.com","password":"password", "role":"user"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "registered successfully")
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("registration with invalid input", func(t *testing.T) {
		// Mock the usecase to return an error when Register is called with invalid data
		mockUserUsecase.On("Register", mock.AnythingOfType("*domain.User")).Return(errors.New("invalid input")).Once()

		w := httptest.NewRecorder()
		reqBody := `{"email":"invalid-email","password":""}` // invalid input
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Expect 500 status due to usecase error
		assert.Contains(t, w.Body.String(), "error")   // Expect an error message
	})
}

func TestUserController_Login(t *testing.T) {
	router := setupRouter()
	mockUserUsecase := new(mocks.UserUsecaseInterface)
	userController := NewUserController(mockUserUsecase)

	router.POST("/login", userController.Login)

	t.Run("successful login", func(t *testing.T) {
		mockUserUsecase.On("Login", mock.AnythingOfType("*domain.User")).Return("token123", nil).Once()

		w := httptest.NewRecorder()
		reqBody := `{"email":"test@example.com","password":"hidden_password"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Debugging output
		t.Logf("Status Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user logged in successfully")
		assert.Contains(t, w.Body.String(), "token")
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("login with invalid input", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"email":"","password":""}` // invalid input
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Debugging output
		t.Logf("Status Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusBadRequest, w.Code) // Expect 400 status
		assert.Contains(t, w.Body.String(), "error")   // Expect an error message
	})

	t.Run("login with incorrect credentials", func(t *testing.T) {
		mockUserUsecase.On("Login", mock.AnythingOfType("*domain.User")).Return("", errors.New("unauthorized access")).Once()

		w := httptest.NewRecorder()
		reqBody := `{"email":"test@example.com","password":"wrongpassword"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Debugging output
		t.Logf("Status Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestUserController_GetUser(t *testing.T) {
	// Initialize the router
	router := setupRouter()

	// Create a new instance of the mock Authorization
	mockAuth := mocks.NewAuthorization(t)

	// Configure the mock to return a gin.HandlerFunc that sets the role
	mockAuth.On("UserAuthorizaiton").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// Example: Set role in the context
			role := c.Query("role") // Adjust according to your role setting method
			c.Set("role", role)
			c.Next()
		}
	}).Once()

	// Use the mock middleware in the router
	router.Use(mockAuth.UserAuthorizaiton())

	// Initialize your controller and routes
	mockUserUsecase := new(mocks.UserUsecaseInterface)
	userController := NewUserController(mockUserUsecase)
	router.GET("/user/:email", userController.GetUser)

	t.Run("admin role gets user", func(t *testing.T) {
		mockUser := &domain.User{Email: "test@example.com", Role: "admin"}
		mockUserUsecase.On("GetUser", "test@example.com").Return(mockUser, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/test@example.com?role=admin", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"email":"test@example.com"`)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("non-admin role tries to get user", func(t *testing.T) {
		// Use the mock middleware to simulate no role set
		mockUser := &domain.User{Email: "test@example.com", Role: "user"}
		mockUserUsecase.On("GetUser", "test@example.com").Return(mockUser, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/test@example.com", nil)
		router.ServeHTTP(w, req)

		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "admin-only")
	})
}

func TestUserController_GetUsers(t *testing.T) {
	// Initialize the router
	router := setupRouter()

	// Create a new instance of the mock Authorization
	mockAuth := mocks.NewAuthorization(t)

	// Configure the mock to return a gin.HandlerFunc that sets the role
	mockAuth.On("UserAuthorizaiton").Return(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// Example: Set role in the context based on query parameter
			role := c.Query("role") // Adjust according to your role setting method
			c.Set("role", role)
			c.Next()
		}
	}).Once()

	// Use the mock middleware in the router
	router.Use(mockAuth.UserAuthorizaiton())

	// Initialize your controller and routes
	mockUserUsecase := new(mocks.UserUsecaseInterface)
	userController := NewUserController(mockUserUsecase)
	router.GET("/users", userController.GetUsers)

	t.Run("admin role gets users", func(t *testing.T) {
		mockUsers := &[]domain.User{
			{Email: "test@example.com", Role: "admin"},
			{Email: "user@example.com", Role: "user"},
		}
		mockUserUsecase.On("GetUsers").Return(mockUsers, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users?role=admin", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"email":"test@example.com"`)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("non-admin role tries to get users", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users?role=user", nil)
		router.ServeHTTP(w, req)

		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "admin-only")
	})
}
