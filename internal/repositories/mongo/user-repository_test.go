package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/dlion/faceit_challenge/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func TestRepository(t *testing.T) {
	t.Run("AddUser", func(t *testing.T) {
		t.Run("Add a new user", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to connect to MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			userRepo := NewUserRepository(mongoClient)
			userRepo.AddUser(ctx, &User{
				FirstName: "testName",
				LastName:  "testLastName",
				Nickname:  "testNickname",
				Email:     "testEmail@email.com",
				Country:   "UK",
				Password:  "testPassword",
			})

			result := mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				FindOne(ctx, bson.M{"nickname": "testNickname", "email": "testEmail@email.com"})

			assert.NoError(t, result.Err())
			userResult := &User{}
			err = result.Decode(userResult)
			assert.NoError(t, err)
			assert.NotNil(t, userResult)
			assert.Equal(t, "testName", userResult.FirstName)
			assert.Equal(t, "testLastName", userResult.LastName)
			assert.Equal(t, "testNickname", userResult.Nickname)
			assert.Equal(t, "testEmail@email.com", userResult.Email)
			assert.Equal(t, "UK", userResult.Country)
			assert.NotEmpty(t, userResult.Id)
			assert.NotEmpty(t, userResult.Password)
			assert.NotEmpty(t, userResult.CreatedAt)
			assert.NotEmpty(t, userResult.UpdatedAt)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte("testPassword")))
		})

		t.Run("Return an error if the user already exist", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			_, err = mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				InsertOne(ctx, &User{
					Id:        "randomID",
					FirstName: "testName",
					LastName:  "testLastName",
					Nickname:  "testNickname",
					Email:     "testEmail@email.com",
					Country:   "UK",
					Password:  "testPwd",
				})
			assert.NoError(t, err)

			userRepo := NewUserRepository(mongoClient)
			err = userRepo.AddUser(ctx, &User{
				FirstName: "testName",
				LastName:  "testLastName",
				Nickname:  "testNickname",
				Email:     "testEmail@email.com",
				Country:   "UK",
				Password:  "testPassword",
			})
			assert.Error(t, err)
		})
	})

	t.Run("Modify an existing user", func(t *testing.T) {

		t.Run("Modify an existing user", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			now := time.Now()
			_, err = mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				InsertOne(ctx, &User{
					Id:        "randomID",
					FirstName: "testName",
					LastName:  "testLastName",
					Nickname:  "testNickname",
					Email:     "testEmail@email.com",
					Country:   "UK",
					Password:  "testPwd",
					CreatedAt: now,
					UpdatedAt: now,
				})
			assert.NoError(t, err)

			userRepo := NewUserRepository(mongoClient)
			err = userRepo.UpdateUser(ctx, &User{
				Id:        "randomID",
				FirstName: "updatedFirstName",
				LastName:  "testLastName",
				Nickname:  "testNickname",
				Email:     "testEmail@email.com",
				Country:   "UK",
				Password:  "testPassword",
			})
			assert.NoError(t, err)

			result := mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				FindOne(ctx, bson.M{"_id": "randomID"})
			assert.NoError(t, result.Err())
			userResult := &User{}
			err = result.Decode(userResult)
			assert.NoError(t, err)
			assert.NotNil(t, userResult)
			assert.Equal(t, "updatedFirstName", userResult.FirstName)
			assert.Equal(t, "testNickname", userResult.Nickname)
			assert.Equal(t, "testEmail@email.com", userResult.Email)
			assert.NotEqual(t, userResult.CreatedAt, userResult.UpdatedAt)
		})

		t.Run("Return an error if the user doesn't exist", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			userRepo := NewUserRepository(mongoClient)
			err = userRepo.UpdateUser(ctx, &User{
				Id:        "randomID",
				FirstName: "updatedFirstName",
				LastName:  "testLastName",
				Nickname:  "testNickname",
				Email:     "testEmail@email.com",
				Country:   "UK",
				Password:  "testPassword",
			})
			assert.Error(t, err)
		})

	})

	t.Run("Remove a user", func(t *testing.T) {
		t.Run("Remove an existing user", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			now := time.Now()
			_, err = mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				InsertOne(ctx, &User{
					Id:        "randomID",
					FirstName: "testName",
					LastName:  "testLastName",
					Nickname:  "testNickname",
					Email:     "testEmail@email.com",
					Country:   "UK",
					Password:  "testPwd",
					CreatedAt: now,
					UpdatedAt: now,
				})
			assert.NoError(t, err)

			userRepo := NewUserRepository(mongoClient)
			err = userRepo.RemoveUser(ctx, &User{Id: "randomID"})

			assert.NoError(t, err)
			result := mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				FindOne(ctx, bson.M{"_id": "randomID"})

			userResult := &User{}
			err = result.Decode(userResult)
			assert.Error(t, err)
		})

		t.Run("Return an error if the user doesn't exist", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			userRepo := NewUserRepository(mongoClient)
			err = userRepo.RemoveUser(ctx, &User{Id: "randomID"})

			assert.Error(t, err)
		})
	})

	t.Run("Return a paginated list of users", func(t *testing.T) {
		t.Run("Just a paginated list of users filtered by country", func(t *testing.T) {
			ctx := context.Background()

			mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
			assert.NoError(t, err, "failed to terminate container: %s", err)

			defer func() {
				err := mongodbContainer.Terminate(ctx)
				assert.NoError(t, err, "failed to terminate container: %s", err)
			}()

			endpoint, err := mongodbContainer.ConnectionString(ctx)
			assert.NoError(t, err, "failed to get connection string: %s", err)

			mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			err = mongoClient.Ping(ctx, nil)
			assert.NoError(t, err, "failed to ping MongoDB: %s", err)

			now := time.Now()
			_, err = mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				InsertOne(ctx, &User{
					Id:        "randomID",
					FirstName: "testName",
					LastName:  "testLastName",
					Nickname:  "testNickname",
					Email:     "testEmail@email.com",
					Country:   "UK",
					Password:  "testPwd",
					CreatedAt: now,
					UpdatedAt: now,
				})
			assert.NoError(t, err)
			_, err = mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				InsertOne(ctx, &User{
					Id:        "randomID1",
					FirstName: "testName",
					LastName:  "testLastName",
					Nickname:  "testNickname1",
					Email:     "testEmail1@email.com",
					Country:   "UK",
					Password:  "testPwd",
					CreatedAt: now,
					UpdatedAt: now,
				})
			assert.NoError(t, err)
			_, err = mongoClient.
				Database(DATABASE_NAME).
				Collection(COLLECTION_NAME).
				InsertOne(ctx, &User{
					Id:        "randomID2",
					FirstName: "testName",
					LastName:  "testLastName",
					Nickname:  "testNickname2",
					Email:     "testEmail2@email.com",
					Country:   "ITA",
					Password:  "testPwd",
					CreatedAt: now,
					UpdatedAt: now,
				})
			assert.NoError(t, err)

			userRepo := NewUserRepository(mongoClient)
			country := "UK"
			userFilter := domain.NewFilterBuilder().ByCountry(&country).Build()
			usersCursor, err := userRepo.GetUsers(ctx, userFilter, int64Ptr(10))
			assert.NoError(t, err)
			defer usersCursor.Close(ctx)

			var users []User
			err = usersCursor.All(ctx, &users)
			assert.NoError(t, err)

			assert.Equal(t, users[0].Id, "randomID")
			assert.Equal(t, users[1].Id, "randomID1")
		})
	})
}
