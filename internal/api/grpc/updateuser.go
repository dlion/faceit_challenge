package grpc

import (
	"context"
	"errors"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	repositories "github.com/dlion/faceit_challenge/internal/repositories/mongo"
	"github.com/dlion/faceit_challenge/pkg/proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserGrpcHandler) UpdateUser(ctx context.Context, request *proto.UpdateUserRequest) (*proto.User, error) {
	serviceReq := user.UpdateUser{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Nickname:  request.GetNickname(),
		Email:     request.GetEmail(),
		Password:  request.GetPassword(),
		Country:   request.GetCountry(),
	}

	user, err := s.userService.UpdateUser(ctx, serviceReq)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, status.Errorf(codes.AlreadyExists, "user not found in the db")
		}

		return nil, status.Error(codes.Internal, "can't update the user")
	}

	return toGrpcUser(user), nil
}
