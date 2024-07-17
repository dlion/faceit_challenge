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
		respUsers[i] = &User{
			Id:        u.Id,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Country:   u.Country,
			Nickname:  u.Nickname,
			CreatedAt: u.CreatedAt.String(),
			UpdatedAt: u.CreatedAt.String(),
		}

	}
	return respUsers, nil
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
