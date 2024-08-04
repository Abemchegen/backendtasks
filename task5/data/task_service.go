package data

import (
	"context"
	"errors"
	"task5/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService struct {
	collection *mongo.Collection
}

func NewTaskService(db *mongo.Database) *TaskService {
	collection := db.Collection("tasks")
	return &TaskService{collection: collection}

}

func (ts *TaskService) CreateTask(newtask *models.Task) (*mongo.InsertOneResult, error) {

	if newtask.Description == "" || newtask.Status == "" || newtask.Title == "" {
		return nil, errors.New("incomplete information")
	}
	return ts.collection.InsertOne(context.TODO(), newtask)

}
func (ts *TaskService) GetTask(id string) (*models.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	var task models.Task

	err = ts.collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&task)

	if err != nil {
		return nil, err
	}

	return &task, nil

}
func (ts *TaskService) GetTasks() (*[]models.Task, error) {
	cursor, err := ts.collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task

	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return nil, err
	}
	return &tasks, nil

}
func (ts *TaskService) UpdateTask(id string, updatedtask *models.Task) (*mongo.UpdateResult, error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	return ts.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": oid},
		bson.D{
			{Key: "$set", Value: updatedtask},
		},
	)

}
func (ts *TaskService) RemoveTask(id string) (*mongo.DeleteResult, error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	return ts.collection.DeleteOne(context.TODO(), bson.M{"_id": oid})

}
