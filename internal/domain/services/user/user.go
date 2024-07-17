package user

import (
	"context"
	"log"

	"github.com/dlion/faceit_challenge/internal/repositories"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(updateUser)
	if err != nil {
		return nil, err
	}

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

func (u *UserServiceImpl) GetUsers(ctx context.Context, query Query) ([]*User, error) {
	log.Printf("Getting users with query: %+v", query)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(query)
	if err != nil {
		return nil, err
	}

	users := u.repository.GetUsers(ctx, NewUserFilterFromQuery(query), query.Limit, query.Offset)

	respUsers := make([]*User, len(users))
	for i, u := range users {
		respUsers[i] = toUser(u)
	}
	return respUsers, nil
}

func toUser(user *repositories.User) *User {
	return &User{
		Id:        user.Id.Hex(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
	}
}

func NewUserFilterFromQuery(query Query) *repositories.UserFilter {
	fbuilder := repositories.NewFilterBuilder()

	if query.FirstName != nil {
		fbuilder = fbuilder.ByFirstName(query.FirstName)
	}

	if query.LastName != nil {
		fbuilder = fbuilder.ByFirstName(query.LastName)
	}

	if query.Nickname != nil {
		fbuilder = fbuilder.ByFirstName(query.Nickname)
	}

	if query.Country != nil {
		fbuilder = fbuilder.ByFirstName(query.Country)
	}

	if query.Email != nil {
		fbuilder = fbuilder.ByFirstName(query.Email)
	}

	return fbuilder.Build()
}
