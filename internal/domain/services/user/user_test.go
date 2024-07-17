package user

import (
	"context"
	"testing"
	"time"

	"github.com/dlion/faceit_challenge/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {
	t.Run("Add a new user and return it", func(t *testing.T) {
		mockedRepository := new(MockUserRepository)
		now := time.Now()
		mockedRepository.On("AddUser").Return(&repositories.User{
			Id:        "66979895733090ace52b13a2",
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "1234567",
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

		userService := NewUserService(mockedRepository)
		addedUser, err := userService.NewUser(context.TODO(), NewUser{
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "testPassword",
		})

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, "66979895733090ace52b13a2", addedUser.Id)
		assert.Equal(t, "TestFirstName", addedUser.FirstName)
		assert.Equal(t, "TestLastName", addedUser.LastName)
		assert.Equal(t, "UK", addedUser.Country)
		assert.Equal(t, "emailTest@test.com", addedUser.Email)
	})
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	args := m.Called()
	return args.Get(0).(*repositories.User), nil
}
