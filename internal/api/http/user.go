package http

import "github.com/dlion/faceit_challenge/internal/domain/services/user"

type UserHandler struct {
	UserService user.UserService
}
