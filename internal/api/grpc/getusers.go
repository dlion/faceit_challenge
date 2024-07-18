package grpc

import (
	"context"

	"github.com/dlion/faceit_challenge/internal"
	"github.com/dlion/faceit_challenge/pkg/proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserGrpcHandler) GetUsers(ctx context.Context, request *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {

	filter := request.GetFilter()
	userFilter := toUserFilter(filter)

	users, err := s.userService.GetUsers(ctx, userFilter)
	if err != nil {
		return nil, status.Error(codes.Internal, "can't get the users")
	}

	userOutput := make([]*proto.User, len(users))
	for i, u := range users {
		userOutput[i] = toGrpcUser(u)
	}

	return &proto.GetUsersResponse{Users: userOutput}, nil
}

func toUserFilter(filter *proto.UserFilter) *internal.UserFilter {
	fbuilder := internal.NewFilterBuilder()

	firstName := filter.FirstName
	if firstName != "" {
		fbuilder.ByFirstName(&firstName)
	}

	lastName := filter.LastName
	if lastName != "" {
		fbuilder = fbuilder.ByLastName(&lastName)
	}

	nickname := filter.Nickname
	if nickname != "" {
		fbuilder = fbuilder.ByNickname(&nickname)
	}

	country := filter.Country
	if country != "" {
		fbuilder = fbuilder.ByCountry(&country)
	}

	email := filter.Email
	if email != "" {
		fbuilder = fbuilder.ByEmail(&email)
	}

	limit := filter.Limit
	if limit > 0 {
		fbuilder.WithLimit(&limit)
	} else {
		fbuilder.WithLimit(intToint64(10))
	}

	offset := filter.Offset
	if offset >= 0 {
		fbuilder.WithOffset(&offset)
	} else {
		fbuilder.WithOffset(intToint64(0))
	}

	return fbuilder.Build()

}

func intToint64(value int) *int64 {
	int64value := int64(value)
	return &int64value
}
