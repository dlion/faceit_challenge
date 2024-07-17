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

	t.Run("Modify and existing user and return it", func(t *testing.T) {
		mockedRepository := new(MockUserRepository)
		now := time.Now()
		later := now.Add(time.Duration(20 * time.Second))
		mockedRepository.On("UpdateUser").Return(&repositories.User{
			Id:        "66979895733090ace52b13a2",
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "1234567",
			CreatedAt: now,
			UpdatedAt: later,
		}, nil)

		userService := NewUserService(mockedRepository)
		updatedUser, err := userService.UpdateUser(context.TODO(), UpdateUser{
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Country:   "UK",
			Email:     "emailTest@test.com",
			Nickname:  "Test",
			Password:  "testPassword",
		})

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, "66979895733090ace52b13a2", updatedUser.Id)
		assert.Equal(t, "TestFirstName", updatedUser.FirstName)
		assert.Equal(t, "TestLastName", updatedUser.LastName)
		assert.Equal(t, "UK", updatedUser.Country)
		assert.Equal(t, "emailTest@test.com", updatedUser.Email)
		assert.Equal(t, later.String(), updatedUser.UpdatedAt)
	})

	t.Run("Remove an existing user", func(t *testing.T) {
		mockedRepository := new(MockUserRepository)
		mockedRepository.On("RemoveUser").Return(nil)

		userService := NewUserService(mockedRepository)
		err := userService.RemoveUser(context.TODO(), "randomId")

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Get paginated list of users filtered by Conuntry", func(t *testing.T) {
		mockedRepository := new(MockUserRepository)
		now := time.Now()
		var dbUsers []*repositories.User
		dbUsers = append(dbUsers, &repositories.User{
			Id:        "66979895733090ace52b13a2",
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
			Id:        "6697c6d1cbfc878bc14673db",
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

		userService := NewUserService(mockedRepository)
		country := "UK"
		users, err := userService.GetUsers(context.TODO(), Query{Country: &country})

		mockedRepository.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "TestFirstName", dbUsers[0].FirstName)
	})

}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	args := m.Called()
	return args.Get(0).(*repositories.User), nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *repositories.User) (*repositories.User, error) {
	args := m.Called()
	return args.Get(0).(*repositories.User), nil
}

func (m *MockUserRepository) RemoveUser(ctx context.Context, id string) error {
	m.Called()
	return nil
}

func (m *MockUserRepository) GetUsers(ctx context.Context, filter repositories.Filter, limit *int64, offset *int64) []*repositories.User {
	args := m.Called()
	return args.Get(0).([]*repositories.User)
}
