package repositories_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"task8/domain"
	"task8/mocks"
	"task8/repositories"
)

func TestRegister(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("successfully registers a user", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "password",
		}

		mockOID := primitive.NewObjectID()
		responseDoc := primitive.E{
			Key:   "insertedId",
			Value: mockOID,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(responseDoc))

		err := repo.Register(mockUser)

		fmt.Printf("Expected ID: %v\n", mockOID)
		fmt.Printf("Actual ID: %v\n", mockUser.ID)

		t.Logf(mockOID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, mockOID, mockUser.ID)
	})

	mt.Run("fails due to MongoDB error", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "password",
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{
				Index:   0,
				Code:    11000,
				Message: "duplicate key error",
			},
		))

		err := repo.Register(mockUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key error")
	})
}

func TestLogin(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully logs in a user", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "password",
		}

		storedUser := domain.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
		}

		// Create a mock for PasswordService
		mockPasswordService := new(mocks.PasswordService)

		// Set expectations for the mock
		mockPasswordService.On("Compare", storedUser.Password, mockUser.Password).Return(nil)

		// Mock MongoDB response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{
			{Key: "email", Value: storedUser.Email},
			{Key: "password", Value: storedUser.Password},
			{Key: "role", Value: storedUser.Role},
		}))

		// Use the mock in the repo.Login
		role, err := repo.Login(mockUser)
		assert.NoError(t, err)
		assert.Equal(t, "user", role)

		// Assert that the expectations were met
		mockPasswordService.AssertExpectations(t)
	})

	mt.Run("fails due to incorrect password", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		storedUser := domain.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
		}

		// Create a mock for PasswordService
		mockPasswordService := new(mocks.PasswordService)

		// Set expectations for the mock
		mockPasswordService.On("Compare", storedUser.Password, mockUser.Password).Return(errors.New("password mismatch"))

		// Mock MongoDB response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{
			{Key: "email", Value: storedUser.Email},
			{Key: "password", Value: storedUser.Password},
			{Key: "role", Value: storedUser.Role},
		}))

		// Use the mock in the repo.Login
		_, err := repo.Login(mockUser)
		assert.Error(t, err)

		// Assert that the expectations were met
		mockPasswordService.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully retrieves a user", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mockUser := &domain.User{
			Email: "test@example.com",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{
			{Key: "email", Value: mockUser.Email},
			{Key: "password", Value: "hashedpassword"},
			{Key: "role", Value: "user"},
		}))

		user, err := repo.GetUser(mockUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, user.Email)
	})

	mt.Run("fails due to user not found", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "users", mtest.FirstBatch))

		user, err := repo.GetUser("nonexistent@example.com")
		assert.Nil(t, user)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})

	mt.Run("fails due to MongoDB error", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "unknown error",
		}))

		_, err := repo.GetUser("test@example.com")
		assert.Error(t, err)
	})
}

func TestGetUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully retrieves users", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mockUsers := []domain.User{
			{
				ID:       primitive.NewObjectID(),
				Email:    "user1@example.com",
				Password: "hashedpassword1",
				Role:     "user",
			},
			{
				ID:       primitive.NewObjectID(),
				Email:    "user2@example.com",
				Password: "hashedpassword2",
				Role:     "admin",
			},
		}

		first := mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockUsers[0].ID},
			{Key: "email", Value: mockUsers[0].Email},
			{Key: "password", Value: mockUsers[0].Password},
			{Key: "role", Value: mockUsers[0].Role},
		})
		second := mtest.CreateCursorResponse(1, "users", mtest.NextBatch, bson.D{
			{Key: "_id", Value: mockUsers[1].ID},
			{Key: "email", Value: mockUsers[1].Email},
			{Key: "password", Value: mockUsers[1].Password},
			{Key: "role", Value: mockUsers[1].Role},
		})
		killCursors := mtest.CreateCursorResponse(0, "users", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		users, err := repo.GetUsers()
		assert.NoError(t, err)
		assert.ElementsMatch(t, mockUsers, *users)
	})

	mt.Run("fails due to MongoDB error", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewUserRepository(mockCollection.Database())

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "unknown error",
		}))

		_, err := repo.GetUsers()
		assert.Error(t, err)
	})
}
