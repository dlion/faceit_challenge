package user

import (
	"context"
	"log"

	"github.com/dlion/faceit_challenge/internal/repositories"
	"github.com/go-playground/validator/v10"
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

	return &User{
		Id:        addedUser.Id,
		FirstName: addedUser.FirstName,
		LastName:  addedUser.LastName,
		Email:     addedUser.Email,
		Country:   addedUser.Country,
		Nickname:  addedUser.Nickname,
		CreatedAt: addedUser.CreatedAt.String(),
		UpdatedAt: addedUser.CreatedAt.String(),
	}, nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, updateUser UpdateUser) (*User, error) {
	log.Printf("Updating user %s", updateUser.Id)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(updateUser)
	if err != nil {
		return nil, err
	}

	repoUser := repositories.NewRepoUser(updateUser.FirstName, updateUser.LastName, updateUser.Nickname, updateUser.Password, updateUser.Email, updateUser.Country)
	repoUser.Id = updateUser.Id
	updatedUser, err := u.repository.UpdateUser(ctx, repoUser)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:        updatedUser.Id,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
		Country:   updatedUser.Country,
		Nickname:  updatedUser.Nickname,
		CreatedAt: updatedUser.CreatedAt.String(),
		UpdatedAt: updatedUser.UpdatedAt.String(),
	}, nil
}
