package repositories

import (
	"context"
	"errors"
	"task7/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepositoryInterface interface {
	CreateTask(newtask *domain.Task, userid string) error
	GetTask(id string) (*domain.Task, error)
	GetTasks(userid string) (*[]domain.Task, error)
	UpdateTask(id string, updatedtask *domain.Task) error
	RemoveTask(id string) error
}

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	collection := db.Collection("tasks")
	return &TaskRepository{collection: collection}

}

func (ts *TaskRepository) CreateTask(newtask *domain.Task, userid string) error {

	userObjectID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return errors.New("user ID is not a valid ObjectID")
	}
	result, err := ts.collection.InsertOne(context.TODO(), newtask)

	if err != nil {
		return errors.New(err.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return errors.New("failed to retrive the inserted ID")
	}

	newtask.ID = oid
	newtask.UserID = userObjectID
	return nil
}

func (ts *TaskRepository) GetTask(id string) (*domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var task domain.Task

	err = ts.collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&task)

	if err != nil {
		return nil, err
	}

	return &task, nil

}
func (ts *TaskRepository) GetTasks(userid string) (*[]domain.Task, error) {
	uid, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return nil, err
	}
	cursor, err := ts.collection.Find(context.TODO(), bson.M{"user_id": uid})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []domain.Task

	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return nil, err
	}
	return &tasks, nil

}
func (ts *TaskRepository) UpdateTask(id string, updatedtask *domain.Task) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	_, err = ts.collection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.D{{Key: "$set", Value: updatedtask}})

	if err != nil {
		return err
	}

	return nil

}
func (ts *TaskRepository) RemoveTask(id string) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = ts.collection.DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		return nil
	}

	return nil

}
