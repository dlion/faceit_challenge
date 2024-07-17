package repositories

import (
	"context"

	"github.com/dlion/faceit_challenge/internal/domain"
)

type UserRepository interface {
	AddUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
	RemoveUser(context.Context, string) error
	GetUsers(context.Context, domain.Filter, *int64, *int64) []*User
}
