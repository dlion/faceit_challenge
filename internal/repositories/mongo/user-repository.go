package repositories

import "go.mongodb.org/mongo-driver/mongo"

type userRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *userRepository {
	return &userRepository{client: client}
}
