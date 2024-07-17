package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddUserHandler(t *testing.T) {
	newUser := user.NewUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Country:   "USA",
	}
	jsonData, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("Failed to marshal new user: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockedUserService := new(MockUserService)
	userHandler := UserHandler{UserService: mockedUserService}
	mockedUserService.On("NewUser").Return(&user.User{
		Id:        "randomCode",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "USA",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	})
	handler := http.HandlerFunc(userHandler.AddUserHandler)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdUser user.User
	err = json.NewDecoder(rr.Body).Decode(&createdUser)
	assert.NoError(t, err)

	assert.NotEmpty(t, createdUser.Id)
	assert.Equal(t, newUser.FirstName, createdUser.FirstName)
	assert.Equal(t, newUser.LastName, createdUser.LastName)
	assert.Equal(t, newUser.Nickname, createdUser.Nickname)
	assert.Equal(t, newUser.Email, createdUser.Email)
	assert.Equal(t, newUser.Country, createdUser.Country)
	assert.NotEmpty(t, createdUser.CreatedAt)
	assert.NotEmpty(t, createdUser.UpdatedAt)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) NewUser(ctx context.Context, newUser user.NewUser) (*user.User, error) {
	args := m.Called()
	return args.Get(0).(*user.User), nil
}
