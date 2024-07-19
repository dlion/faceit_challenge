package grpc

import (
	"context"
	"errors"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	repositories "github.com/dlion/faceit_challenge/internal/repositories/mongo"
	"github.com/dlion/faceit_challenge/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGrpcHandler struct {
	proto.UnimplementedUserServiceServer
	userService user.UserService
}

func NewUserGrpcHandler(userService user.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{
		userService: userService,
	}
}

func (s *UserGrpcHandler) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.User, error) {
	serviceReq := &user.NewUser{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Nickname:  request.GetNickname(),
		Email:     request.GetEmail(),
		Password:  request.GetPassword(),
		Country:   request.GetCountry(),
	}

	user, err := s.userService.NewUser(ctx, serviceReq)
	if err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExist) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exist in the db")
		}

		return nil, status.Error(codes.Internal, "can't create the user")
	}

	return toGrpcUser(user), nil
}

func toGrpcUser(user *user.User) *proto.User {
	return &proto.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
