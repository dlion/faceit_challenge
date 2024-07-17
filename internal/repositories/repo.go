package repositories

import (
	"context"

	"github.com/dlion/faceit_challenge/internal"
)

type UserRepository interface {
	AddUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
	RemoveUser(context.Context, string) error
	GetUsers(context.Context, *internal.UserFilter, *int64, *int64) ([]*User, error)
}
