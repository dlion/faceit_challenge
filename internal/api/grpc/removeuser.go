package grpc

import (
	"context"
	"errors"

	repositories "github.com/dlion/faceit_challenge/internal/repositories/mongo"
	"github.com/dlion/faceit_challenge/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserGrpcHandler) RemoveUser(ctx context.Context, request *proto.DeleteUserRequest) error {

	err := s.userService.RemoveUser(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return status.Errorf(codes.AlreadyExists, "user not found in the db")
		}

		return status.Error(codes.Internal, "can't remove the user")
	}

	return nil
}
