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
