package handlers

import "github.com/dlion/faceit_challenge/internal/domain/services/user"

type UserHandler struct {
	UserService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}
