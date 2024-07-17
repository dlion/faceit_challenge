package repositories

import (
	"context"
)

type UserRepository interface {
	AddUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
}
