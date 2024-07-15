package repositories

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE_NAME   = "faceit"
	COLLECTION_NAME = "users"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{collection: client.Database("faceit").Collection("users")}
}

func (u *UserRepository) AddUser(ctx context.Context, user User) error {
	log.Println("Adding user to the database")
	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
