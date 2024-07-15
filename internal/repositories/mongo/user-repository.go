package repositories

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	DATABASE_NAME   = "faceit"
	COLLECTION_NAME = "users"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{collection: client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)}
}

func (u *UserRepository) AddUser(ctx context.Context, user *User) error {
	log.Println("Adding a user to the database")

	err := addHashedPassword(user)
	if err != nil {
		return err
	}

	_, err = u.collection.InsertOne(ctx, user)
	if err != nil {
		return err
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
