package repositories_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"golang.org/x/crypto/bcrypt"

	"task8/domain"
	"task8/mocks"
	"task8/repositories"
)

func TestRegister(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully registers a user", func(mt *mtest.T) {

		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "password",
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mockps.On("Hash", mock.Anything).Return("hashedpassword", nil)
		err := repo.Register(mockUser)

		t.Logf("Inserted OID: %s", mockUser.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, mockUser.ID)
	})

	mt.Run("fails due to MongoDB error", func(mt *mtest.T) {
		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "password",
		}
		mockps.On("Hash", mock.Anything).Return("hashedpassword", nil)

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

	// Generate a bcrypt hash for the password
	plainPassword := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	mt.Run("successfully logs in a user", func(mt *mtest.T) {
		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: plainPassword, // Plaintext password
		}
		storedUser := domain.User{
			Email:    "test@example.com",
			Password: string(hashedPassword), // Use the hashed password
			Role:     "user",
		}

		mockps.On("Compare", mock.Anything, mock.Anything).Return(nil)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{Key: "email", Value: storedUser.Email},
			{Key: "password", Value: storedUser.Password},
			{Key: "role", Value: storedUser.Role},
		}))

		// Use the mock in the repo.Login
		role, err := repo.Login(mockUser)

		if err != nil {
			t.Logf("Error: %v", err)
		}
		assert.NoError(t, err)
		assert.Equal(t, "user", role)

		mockps.AssertExpectations(t)
	})

	mt.Run("fails due to incorrect password", func(mt *mtest.T) {
		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mockUser := &domain.User{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		storedUser := domain.User{
			Email:    "test@example.com",
			Password: string(hashedPassword), // Use the hashed password
			Role:     "user",
		}

		mockps.On("Compare", mock.Anything, mock.Anything).Return(errors.New("passwords do not match"))

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{Key: "email", Value: storedUser.Email},
			{Key: "password", Value: storedUser.Password},
			{Key: "role", Value: storedUser.Role},
		}))

		_, err := repo.Login(mockUser)
		assert.Error(t, err)

		mockps.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully retrieves a user", func(mt *mtest.T) {
		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mockUser := &domain.User{
			Email: "test@example.com",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
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
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		user, err := repo.GetUser("nonexistent@example.com")
		assert.Nil(t, user)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})

	mt.Run("fails due to MongoDB error", func(mt *mtest.T) {
		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

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
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

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

		first := mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockUsers[0].ID},
			{Key: "email", Value: mockUsers[0].Email},
			{Key: "password", Value: mockUsers[0].Password},
			{Key: "role", Value: mockUsers[0].Role},
		})
		second := mtest.CreateCursorResponse(1, "test.users", mtest.NextBatch, bson.D{
			{Key: "_id", Value: mockUsers[1].ID},
			{Key: "email", Value: mockUsers[1].Email},
			{Key: "password", Value: mockUsers[1].Password},
			{Key: "role", Value: mockUsers[1].Role},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.users", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		users, err := repo.GetUsers()
		assert.NoError(t, err)
		assert.ElementsMatch(t, mockUsers, *users)
	})

	mt.Run("fails due to MongoDB error", func(mt *mtest.T) {
		mockCollection := mt.Coll
		mockps := new(mocks.PasswordService)
		repo := repositories.NewUserRepository(mockCollection.Database(), mockps)

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "unknown error",
		}))

		_, err := repo.GetUsers()
		assert.Error(t, err)
	})
}
