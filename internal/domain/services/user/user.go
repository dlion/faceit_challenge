package user

import (
	"context"
	"log"
	"time"

	"github.com/dlion/faceit_challenge/internal"
	"github.com/dlion/faceit_challenge/internal/repositories"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	NewUser(context.Context, NewUser) (*User, error)
	UpdateUser(context.Context, UpdateUser) (*User, error)
	RemoveUser(context.Context, string) error
	GetUsers(context.Context, *internal.UserFilter) ([]*User, error)
}

type UserServiceImpl struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repository: repository}
}

func (u *UserServiceImpl) NewUser(ctx context.Context, newUser NewUser) (*User, error) {
	log.Printf("Adding a new user: %+v", newUser)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(newUser)
	if err != nil {
		return nil, err
	}

	repoUser := repositories.NewRepoUser(newUser.FirstName, newUser.LastName, newUser.Nickname, newUser.Password, newUser.Email, newUser.Country)
	addedUser, err := u.repository.AddUser(ctx, repoUser)
	if err != nil {
		return nil, err
	}

	return toUser(addedUser), nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, updateUser UpdateUser) (*User, error) {
	log.Printf("Updating user %s", updateUser.Id)

	repoUser := repositories.NewRepoUser(updateUser.FirstName, updateUser.LastName, updateUser.Nickname, updateUser.Password, updateUser.Email, updateUser.Country)

	hex, err := primitive.ObjectIDFromHex(updateUser.Id)
	if err != nil {
		return nil, err
	}
	repoUser.Id = hex

	updatedUser, err := u.repository.UpdateUser(ctx, repoUser)
	if err != nil {
		return nil, err
	}

	return toUser(updatedUser), nil
}

func (u *UserServiceImpl) RemoveUser(ctx context.Context, id string) error {
	log.Printf("Removing user with id: %s", id)

	err := u.repository.RemoveUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) GetUsers(ctx context.Context, userFilter *internal.UserFilter) ([]*User, error) {
	log.Printf("Getting users with query: %+v", userFilter)

	users, err := u.repository.GetUsers(ctx, userFilter, userFilter.Limit, userFilter.Offset)
	if err != nil {
		return nil, err
	}

	respUsers := make([]*User, len(users))
	for i, u := range users {
		respUsers[i] = toUser(u)
	}
	return respUsers, nil
}

func toUser(user *repositories.User) *User {
	return &User{
		Id:        user.Id.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
