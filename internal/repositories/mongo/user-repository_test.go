package repositories

import (
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestRepository(t *testing.T) {
	t.Run("Add a new user", func(t *testing.T) {
		ctx := context.Background()

		mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
		if err != nil {
			log.Fatalf("failed to start container: %s", err)
		}

		defer func() {
			if err := mongodbContainer.Terminate(ctx); err != nil {
				log.Fatalf("failed to terminate container: %s", err)
			}
		}()

		endpoint, err := mongodbContainer.ConnectionString(ctx)
		if err != nil {
			log.Fatalf("failed to get connection string: %s", err)
		}

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
		if err != nil {
			log.Fatalf("failed to connect to MongoDB: %s", err)
		}

		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("failed to ping MongoDB: %s", err)
		}

		//userRepo := NewUserRepository(mongoClient)
		// userRepo.AddUser()

		// result := mongoClient.
		// 	Database("faceit").
		// 	Collection("users").
		// 	FindOne(ctx, bson.M{"nickname": "testNickname", "email": "testEmail@email.com"})

		// assert.NoError(t, result.Err())
		// userResult := &User{}
		// assert.Contains(t, result.Decode(userResult), &User{
		// 	FirstName: "testName",
		// 	LastName:  "testLastName",
		// 	Nickname:  "testNickname",
		// 	Email:     "testEmail@email.com",
		// 	Country:   "UK",
		// })
	})

	t.Run("Modify an existing user", func(t *testing.T) {

	})

	t.Run("Remove a user", func(t *testing.T) {

	})

	t.Run("Return a paginated list of users, filtering by some criterias", func(t *testing.T) {

	})
}
