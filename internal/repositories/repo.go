package repositories

import (
	"context"
)

type UserRepository interface {
	AddUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
	RemoveUser(context.Context, string) error
	GetUsers(context.Context, Filter, *int64, *int64) []*User
}
