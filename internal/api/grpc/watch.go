package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/dlion/faceit_challenge/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserGrpcHandler) Watch(_ *emptypb.Empty, server proto.UserService_WatchServer) error {
	clientId := uuid.New().String()
	channel := s.userService.GetChangeChannel(clientId)
	defer s.userService.RemoveChannel(clientId)

	for {
		select {
		// Check if the client has disconnected or the context has been canceled
		case <-server.Context().Done():
			err := server.Context().Err()

			if errors.Is(err, context.Canceled) {
				log.Print("Client disconnected")
			} else {
				log.Print("Client disconnected with error, ", err)
			}

			return nil
		case change, closed := <-channel:
			if !closed {
				return nil
			}

			response := &proto.WatchResponse{
				ChangeType: change.OperationType,
				UserId:     change.UserId,
			}

			err := server.Send(response)
			if err != nil {
				return err
			}
		}
	}
}
