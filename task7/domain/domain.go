package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty" json:"-"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	DueDate     time.Time          `json:"duedate"`
	Status      string             `json:"status"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Role     string             `bson:"role" json:"role"`
}
type TaskRepositoryInterface interface {
	CreateTask(newtask *Task, userid string) error
	GetTask(id string) (*Task, error)
	GetTasks(userid string) (*[]Task, error)
	UpdateTask(id string, updatedtask *Task) error
	RemoveTask(id string) error
}
type UserRepositoryInterface interface {
	Register(user *User) error
	Login(user *User) (string, error)
	GetUser(email string) (*User, error)
	GetUsers() (*[]User, error)
}
