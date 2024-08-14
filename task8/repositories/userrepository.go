package repositories

import (
	"context"
	"errors"
	"task8/domain"
	"task8/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
	password   infrastructure.PasswordService
}

func NewUserRepository(db *mongo.Database, ps infrastructure.PasswordService) *UserRepository {
	collection := db.Collection("users")
	return &UserRepository{collection: collection, password: ps}
}

func (us *UserRepository) Register(user *domain.User) error {

	hashedPassword, err := us.password.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

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

func (us *UserRepository) Login(user *domain.User) (string, error) {

	var u domain.User
	err := us.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u)

	if err != nil {
		return u.Role, err
	}

	err = us.password.Compare(u.Password, user.Password)

	if err != nil {
		return u.Role, err
	}
	return u.Role, nil

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
