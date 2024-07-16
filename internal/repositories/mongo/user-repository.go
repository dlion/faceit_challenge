package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dlion/faceit_challenge/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	DATABASE_NAME   = "faceit"
	COLLECTION_NAME = "users"
)

var (
	ErrUserAlreadyExist = errors.New("the user already exist in the db")
	ErrUserNotFound     = errors.New("the user doesn't exist in the db")
	ErrNothingToUpdate  = errors.New("there's anything to update")
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{collection: client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)}
}

func (u *UserRepository) AddUser(ctx context.Context, user *User) error {
	log.Printf("Adding a user to the database")

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

func (u *UserRepository) UpdateUser(ctx context.Context, user *User) error {
	log.Printf("Updating user (%s) in the database", user.Id)

	updatedUser, err := createUpdatedUser(user)
	if err != nil {
		return err
	}

	updatedResult, err := u.collection.UpdateByID(ctx, user.Id, updatedUser)
	if err != nil {
		return err
	}

	if updatedResult.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (u *UserRepository) RemoveUser(ctx context.Context, user *User) error {
	log.Printf("Removing user (%s) from the database", user.Id)

	deletedResult, err := u.collection.DeleteOne(ctx, bson.M{"_id": user.Id})
	if err != nil {
		return err
	}

	if deletedResult.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (u *UserRepository) GetUsers(ctx context.Context, userFilter domain.Filter, limit *int64) (*mongo.Cursor, error) {

	log.Printf("Getting users from the database with filters: %+v", userFilter.ToBSON())

	if limit == nil {
		limit = int64Ptr(10)
	}

	cursor, err := u.collection.Find(ctx, userFilter.ToBSON(), &options.FindOptions{
		Limit: limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func int64Ptr(value int64) *int64 {
	return &value
}

func createUpdatedUser(user *User) (bson.M, error) {
	updateFields := bson.M{}

	if user.FirstName != "" {
		updateFields["first_name"] = user.FirstName
	}

	if user.LastName != "" {
		updateFields["last_name"] = user.LastName
	}

	if user.Nickname != "" {
		updateFields["nickname"] = user.Nickname
	}

	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		updateFields["password"] = hashedPassword
	}

	if user.Email != "" {
		updateFields["email"] = user.Email
	}

	if user.Country != "" {
		updateFields["country"] = user.Country
	}

	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 0 {
		return nil, ErrNothingToUpdate
	}

	return bson.M{"$set": updateFields}, nil
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
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func setCreationTime(user *User) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
}
