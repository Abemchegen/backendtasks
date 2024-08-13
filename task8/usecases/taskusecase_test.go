package usecases_test

import (
	"task8/domain"
	"task8/mocks"
	"task8/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryInterface)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	task := &domain.Task{
		Title:       "Sample Task",
		Description: "This is a sample task",
		Status:      "pending",
	}

	mockRepo.On("CreateTask", mock.AnythingOfType("*domain.Task"), "userID").Return(nil)

	err := taskUsecase.CreateTask(task, "userID")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTask(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryInterface)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	taskID := primitive.NewObjectID()
	task := &domain.Task{
		ID:          taskID,
		Title:       "Sample Task",
		Description: "This is a sample task",
		Status:      "pending",
	}

	mockRepo.On("GetTask", taskID.Hex()).Return(task, nil)

	result, err := taskUsecase.GetTask(taskID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, task, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTasks(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryInterface)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	userID := primitive.NewObjectID()
	tasks := &[]domain.Task{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Title:       "Task 1",
			Description: "This is task 1",
			Status:      "pending",
		},
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			Title:       "Task 2",
			Description: "This is task 2",
			Status:      "completed",
		},
	}

	mockRepo.On("GetTasks", userID.Hex()).Return(tasks, nil)

	result, err := taskUsecase.GetTasks(userID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, tasks, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryInterface)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	taskID := primitive.NewObjectID()
	updatedTask := &domain.Task{
		Title:       "Updated Task",
		Description: "This is the updated task",
		Status:      "in-progress",
	}

	mockRepo.On("UpdateTask", taskID.Hex(), updatedTask).Return(nil)

	err := taskUsecase.UpdateTask(taskID.Hex(), updatedTask)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTask(t *testing.T) {
	mockRepo := new(mocks.TaskRepositoryInterface)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	taskID := primitive.NewObjectID()

	mockRepo.On("RemoveTask", taskID.Hex()).Return(nil)

	err := taskUsecase.RemoveTask(taskID.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
