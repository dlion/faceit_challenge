package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	DATABASE_NAME   = "faceit"
	COLLECTION_NAME = "users"
)

var (
	ErrUserAlreadyExist = errors.New("the user already exist in the db")
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{collection: client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)}
}

func (u *UserRepository) AddUser(ctx context.Context, user *User) error {
	log.Println("Adding a user to the database")

	if err := userAlreadyExists(ctx, u.collection, user.Nickname, user.Email); err != nil {
		return err
	}

	if err := addHashedPassword(user); err != nil {
		return err
	}

	setCreationTime(user)

	if _, err := u.collection.InsertOne(ctx, user); err != nil {
		return err
	}

	return nil
}

func userAlreadyExists(ctx context.Context, collection *mongo.Collection, nickname, email string) error {
	count, err := collection.CountDocuments(ctx, bson.M{"nickname": nickname, "email": email})
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrUserAlreadyExist
	}

	return nil
}

func addHashedPassword(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

func setCreationTime(user *User) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
}
