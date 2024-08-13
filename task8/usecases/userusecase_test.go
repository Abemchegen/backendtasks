package usecases

import (
	"errors"
	"task8/domain"
	"task8/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRegister(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	userUsecase := domain.UserUsecaseInterface(NewUserUsecase(mockRepo))

	user := &domain.User{
		Email:    "test@example.com",
		Password: "password",
		Role:     "user",
	}

	mockRepo.On("Register", user).Return(nil)

	err := userUsecase.Register(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterIncompleteInfo(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	userUsecase := domain.UserUsecaseInterface(NewUserUsecase(mockRepo))

	user := &domain.User{
		Email:    "",
		Password: "password",
		Role:     "",
	}

	err := userUsecase.Register(user)

	assert.Error(t, err)
	assert.Equal(t, "incomplete information", err.Error())
}

func TestLogin(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	userUsecase := domain.UserUsecaseInterface(NewUserUsecase(mockRepo))

	user := &domain.User{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.On("Login", user).Return("user", nil)

	token, err := userUsecase.Login(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	userUsecase := domain.UserUsecaseInterface(NewUserUsecase(mockRepo))

	user := &domain.User{
		Email:    "invalid@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("Login", user).Return("", errors.New("invalid email or password"))

	token, err := userUsecase.Login(user)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	userUsecase := domain.UserUsecaseInterface(NewUserUsecase(mockRepo))

	email := "test@example.com"
	expectedUser := &domain.User{
		ID:    primitive.NewObjectID(),
		Email: email,
		Role:  "user",
	}

	mockRepo.On("GetUser", email).Return(expectedUser, nil)

	user, err := userUsecase.GetUser(email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)
	userUsecase := domain.UserUsecaseInterface(NewUserUsecase(mockRepo))

	expectedUsers := &[]domain.User{
		{ID: primitive.NewObjectID(), Email: "user1@example.com", Role: "user"},
		{ID: primitive.NewObjectID(), Email: "user2@example.com", Role: "user"},
	}

	mockRepo.On("GetUsers").Return(expectedUsers, nil)

	users, err := userUsecase.GetUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}
