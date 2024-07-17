package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dlion/faceit_challenge/internal/domain"
	"github.com/dlion/faceit_challenge/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type UserRepositoryMongoImpl struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongoImpl(client *mongo.Client) *UserRepositoryMongoImpl {
	return &UserRepositoryMongoImpl{collection: client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)}
}

func (u *UserRepositoryMongoImpl) AddUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	log.Printf("Adding a user to the database")

	if err := userAlreadyExists(ctx, u.collection, user.Nickname, user.Email); err != nil {
		return nil, err
	}

	if err := addHashedPassword(user); err != nil {
		return nil, err
	}

	setCreationTime(user)

	insertedUserID, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	insertedObjectID, ok := insertedUserID.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatalf("Failed to convert the insertedID to an ObjectID")
	}

	user.Id = insertedObjectID.Hex()

	return user, nil
}

func (u *UserRepositoryMongoImpl) UpdateUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	log.Printf("Updating user (%s) in the database", user.Id)

	updatedUser, err := createUpdatedUser(user)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return nil, err
	}

	updatedResult, err := u.collection.UpdateByID(ctx, userId.String(), updatedUser)
	if err != nil {
		return nil, err
	}

	if updatedResult.MatchedCount == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (u *UserRepositoryMongoImpl) RemoveUser(ctx context.Context, id string) error {
	log.Printf("Removing user (%s) from the database", id)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	deletedResult, err := u.collection.DeleteOne(ctx, bson.M{"_id": objectId.String()})
	if err != nil {
		return err
	}

	if deletedResult.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (u *UserRepositoryMongoImpl) GetUsers(ctx context.Context, userFilter domain.Filter, limit, offset *int64) (*mongo.Cursor, error) {

	log.Printf("Getting users from the database with filters: %+v", userFilter.ToBSON())

	if limit == nil {
		limit = int64Ptr(10)
	}

	if offset == nil {
		offset = int64Ptr(0)
	}

	cursor, err := u.collection.Find(ctx, userFilter.ToBSON(), &options.FindOptions{
		Limit: limit,
		Skip:  offset,
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

func createUpdatedUser(user *repositories.User) (bson.M, error) {
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

func addHashedPassword(user *repositories.User) error {
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

func setCreationTime(user *repositories.User) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
}
