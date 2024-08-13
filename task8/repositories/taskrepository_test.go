package repositories_test

import (
	"task8/domain"
	"task8/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully creates a task", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		mockTask := &domain.Task{
			Title:       "Sample Task",
			Description: "This is a sample task",
			Status:      "pending",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.CreateTask(mockTask, primitive.NewObjectID().Hex())

		assert.NoError(t, err)
		assert.NotNil(t, mockTask.ID)
	})

	mt.Run("fails due to invalid userID", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		mockTask := &domain.Task{
			Title:       "Sample Task",
			Description: "This is a sample task",
			Status:      "pending",
		}

		err := repo.CreateTask(mockTask, "invalidUserID")

		assert.Error(t, err)
		assert.EqualError(t, err, "user ID is not a valid ObjectID")
	})
}

func TestGetTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully retrieves a task", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		taskID := primitive.NewObjectID()
		mockTask := &domain.Task{
			ID:          taskID,
			Title:       "Sample Task",
			Description: "This is a sample task",
			Status:      "pending",
		}

		first := mtest.CreateCursorResponse(1, "tasks.task", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: taskID},
			{Key: "title", Value: "Sample Task"},
			{Key: "description", Value: "This is a sample task"},
			{Key: "status", Value: "pending"},
		})
		killCursors := mtest.CreateCursorResponse(0, "tasks.task", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		task, err := repo.GetTask(taskID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, mockTask, task)
	})

	mt.Run("fails due to invalid ObjectID", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		_, err := repo.GetTask("invalidID")

		assert.Error(t, err)
		assert.EqualError(t, err, "the provided hex string is not a valid ObjectID")
	})

	mt.Run("fails due to task not found", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		taskID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "tasks.task", mtest.FirstBatch))

		task, err := repo.GetTask(taskID.Hex())

		assert.Nil(t, task)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})
}

func TestGetTasks(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully retrieves tasks", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		userID := primitive.NewObjectID()
		mockTasks := []domain.Task{
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

		first := mtest.CreateCursorResponse(1, "tasks.task", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockTasks[0].ID},
			{Key: "user_id", Value: userID},
			{Key: "title", Value: "Task 1"},
			{Key: "description", Value: "This is task 1"},
			{Key: "status", Value: "pending"},
		})
		second := mtest.CreateCursorResponse(1, "tasks.task", mtest.NextBatch, bson.D{
			{Key: "_id", Value: mockTasks[1].ID},
			{Key: "user_id", Value: userID},
			{Key: "title", Value: "Task 2"},
			{Key: "description", Value: "This is task 2"},
			{Key: "status", Value: "completed"},
		})
		killCursors := mtest.CreateCursorResponse(0, "tasks.task", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		tasks, err := repo.GetTasks(userID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, &mockTasks, tasks)
	})

	mt.Run("fails due to invalid userID", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		_, err := repo.GetTasks("invalidUserID")

		assert.Error(t, err)
		assert.EqualError(t, err, "the provided hex string is not a valid ObjectID")
	})
}

func TestUpdateTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully updates a task", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		taskID := primitive.NewObjectID()
		updatedTask := &domain.Task{
			Title:       "Updated Task",
			Description: "This is the updated task",
			Status:      "in-progress",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.UpdateTask(taskID.Hex(), updatedTask)

		assert.NoError(t, err)
	})

	mt.Run("fails due to invalid ObjectID", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		updatedTask := &domain.Task{
			Title:       "Updated Task",
			Description: "This is the updated task",
			Status:      "in-progress",
		}

		err := repo.UpdateTask("invalidID", updatedTask)

		assert.Error(t, err)
		assert.EqualError(t, err, "the provided hex string is not a valid ObjectID")
	})
}

func TestRemoveTask(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successfully removes a task", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		taskID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.RemoveTask(taskID.Hex())

		assert.NoError(t, err)
	})

	mt.Run("fails due to invalid ObjectID", func(mt *mtest.T) {
		mockCollection := mt.Coll
		repo := repositories.NewTaskRepository(mockCollection.Database())

		err := repo.RemoveTask("invalidID")

		assert.Error(t, err)
		assert.EqualError(t, err, "the provided hex string is not a valid ObjectID")
	})
}
