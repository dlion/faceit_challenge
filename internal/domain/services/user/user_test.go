package user

import (
	"context"
	"testing"
	"time"

	filter "github.com/dlion/faceit_challenge/internal"
	"github.com/dlion/faceit_challenge/internal/repositories"
	"github.com/dlion/faceit_challenge/pkg/notifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserService(t *testing.T) {
	t.Run("Add a new user and return it", func(t *testing.T) {
		mockedRepository := new(mockUserRepository)
		mockedNotifier := new(mockUserNotifier)
		now := time.Now()
		objectId := primitive.NewObjectIDFromTimestamp(now)
		mockedRepository.On("AddUser").Return(&repositories.User{
			Id:        objectId,
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "1234567",
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)
		mockedNotifier.On("Broadcast")

		userService := NewUserService(mockedRepository, mockedNotifier)
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
		assert.Equal(t, objectId.Hex(), addedUser.Id)
		assert.Equal(t, "TestFirstName", addedUser.FirstName)
		assert.Equal(t, "TestLastName", addedUser.LastName)
		assert.Equal(t, "UK", addedUser.Country)
		assert.Equal(t, "emailTest@test.com", addedUser.Email)
	})

	t.Run("Modify and existing user and return it", func(t *testing.T) {
		mockedRepository := new(mockUserRepository)
		mockedNotifier := new(mockUserNotifier)
		now := time.Now()
		objectId := primitive.NewObjectIDFromTimestamp(now)
		later := now.Add(time.Duration(20 * time.Second))
		mockedRepository.On("UpdateUser").Return(&repositories.User{
			Id:        objectId,
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "1234567",
			CreatedAt: now,
			UpdatedAt: later,
		}, nil)
		mockedNotifier.On("Broadcast")

		userService := NewUserService(mockedRepository, mockedNotifier)
		updatedUser, err := userService.UpdateUser(context.TODO(), UpdateUser{
			Id:        objectId.Hex(),
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "testPassword",
		})

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, objectId.Hex(), updatedUser.Id)
		assert.Equal(t, "TestFirstName", updatedUser.FirstName)
		assert.Equal(t, "TestLastName", updatedUser.LastName)
		assert.Equal(t, "UK", updatedUser.Country)
		assert.Equal(t, "emailTest@test.com", updatedUser.Email)
		assert.Equal(t, later.Format(time.RFC3339), updatedUser.UpdatedAt)
	})

	t.Run("Remove an existing user", func(t *testing.T) {
		mockedRepository := new(mockUserRepository)
		mockedNotifier := new(mockUserNotifier)
		mockedRepository.On("RemoveUser").Return(nil)
		mockedNotifier.On("Broadcast")

		userService := NewUserService(mockedRepository, mockedNotifier)
		err := userService.RemoveUser(context.TODO(), "randomId")

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Get paginated list of users filtered by Conuntry", func(t *testing.T) {
		mockedRepository := new(mockUserRepository)
		mockedNotifier := new(mockUserNotifier)
		now := time.Now()
		objectId := primitive.NewObjectIDFromTimestamp(now)
		var dbUsers []*repositories.User
		dbUsers = append(dbUsers, &repositories.User{
			Id:        objectId,
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "1234567",
			CreatedAt: now,
			UpdatedAt: now,
		})

		later := now.Add(time.Duration(20 * time.Second))
		dbUsers = append(dbUsers, &repositories.User{
			Id:        objectId,
			FirstName: "TestFirstName1",
			LastName:  "TestLastName1",
			Country:   "UK",
			Email:     "emailTest1@test.com",
			Nickname:  "Test1",
			Password:  "12345678",
			CreatedAt: later,
			UpdatedAt: later,
		})
		mockedRepository.On("GetUsers").Return(dbUsers, nil)

		userService := NewUserService(mockedRepository, mockedNotifier)
		country := "UK"
		users, err := userService.GetUsers(context.TODO(), &filter.UserFilter{Country: &country})

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "TestFirstName", dbUsers[0].FirstName)
	})

}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) AddUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	args := m.Called()
	return args.Get(0).(*repositories.User), nil
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	args := m.Called()
	return args.Get(0).(*repositories.User), nil
}

func (m *mockUserRepository) RemoveUser(ctx context.Context, id string) error {
	m.Called()
	return nil
}

func (m *mockUserRepository) GetUsers(ctx context.Context, filter *filter.UserFilter, limit *int64, offset *int64) ([]*repositories.User, error) {
	args := m.Called()
	return args.Get(0).([]*repositories.User), nil
}

type mockUserNotifier struct {
	mock.Mock
}

func (m *mockUserNotifier) AddSubscriber(id string) <-chan notifier.ChangeData {
	m.Called()
	return nil
}

func (m *mockUserNotifier) RemoveSubscriber(id string) {
	m.Called()
}

func (m *mockUserNotifier) Broadcast(msg notifier.ChangeData) {
	m.Called()
}

func (m *mockUserNotifier) Close() {
	m.Called()
}
