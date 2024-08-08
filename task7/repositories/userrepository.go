package repositories

import (
	"context"
	"errors"
	"task7/domain"
	"task7/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection("users")
	return &UserRepository{collection: collection}
}

func (us *UserRepository) Register(user *domain.User) error {

	result, err := us.collection.InsertOne(context.TODO(), user)

	if err != nil {
		return err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return errors.New("failed to retrive the inserted ID")
	}
	user.ID = oid
	return nil

}

func (us *UserRepository) Login(user *domain.User) error {

	var u domain.User
	err := us.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u)

	if err != nil {
		return errors.New("invalid email or password")
	}

	err = infrastructure.Compare(u.Password, user.Password)

	if err != nil {
		return errors.New("invalid email or password")
	}
	return nil

}

func (us *UserRepository) GetUser(email string) (*domain.User, error) {

	var user domain.User

	err := us.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserRepository) GetUsers() (*[]domain.User, error) {

	cursor, err := us.collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []domain.User

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return &users, nil

}
