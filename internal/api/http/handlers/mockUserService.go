package handlers

import (
	"context"

	filter "github.com/dlion/faceit_challenge/internal"
	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	"github.com/dlion/faceit_challenge/pkg/notifier"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) NewUser(ctx context.Context, newUser user.NewUser) (*user.User, error) {
	args := m.Called()
	return args.Get(0).(*user.User), nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, updateUser user.UpdateUser) (*user.User, error) {
	args := m.Called()
	return args.Get(0).(*user.User), nil
}

func (m *MockUserService) RemoveUser(ctx context.Context, id string) error {
	m.Called()
	return nil
}

func (m *MockUserService) GetUsers(ctx context.Context, filter *filter.UserFilter) ([]*user.User, error) {
	args := m.Called()
	return args.Get(0).([]*user.User), nil
}

func (m *MockUserService) GetChangeChannel(clientId string) <-chan notifier.ChangeData {
	args := m.Called()
	return args.Get(0).(<-chan notifier.ChangeData)
}

func (m *MockUserService) RemoveChannel(clientId string) error {
	m.Called()
	return nil
}
