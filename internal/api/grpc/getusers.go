package grpc

import (
	"context"

	filter "github.com/dlion/faceit_challenge/internal"
	"github.com/dlion/faceit_challenge/pkg/proto"
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

func toUserFilter(userFilter *proto.UserFilter) *filter.UserFilter {
	fbuilder := filter.NewFilterBuilder()

	firstName := userFilter.FirstName
	if firstName != "" {
		fbuilder.ByFirstName(&firstName)
	}

	lastName := userFilter.LastName
	if lastName != "" {
		fbuilder = fbuilder.ByLastName(&lastName)
	}

	nickname := userFilter.Nickname
	if nickname != "" {
		fbuilder = fbuilder.ByNickname(&nickname)
	}

	country := userFilter.Country
	if country != "" {
		fbuilder = fbuilder.ByCountry(&country)
	}

	email := userFilter.Email
	if email != "" {
		fbuilder = fbuilder.ByEmail(&email)
	}

	limit := userFilter.Limit
	if limit > 0 {
		fbuilder.WithLimit(&limit)
	} else {
		fbuilder.WithLimit(intToint64(10))
	}

	offset := userFilter.Offset
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
