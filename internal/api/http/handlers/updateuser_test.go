package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserHandler(t *testing.T) {
	updateUser := user.UpdateUser{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Country:   "USA",
	}
	jsonData, err := json.Marshal(updateUser)
	if err != nil {
		t.Fatalf("Failed to marshal update user: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/user/66981a71a4fd0f7ff33251b1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockedUserService := new(MockUserService)
	userHandler := UserHandler{UserService: mockedUserService}
	mockedUserService.On("UpdateUser").Return(&user.User{
		Id:        "66981a71a4fd0f7ff33251b1",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "USA",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	})
	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", userHandler.UpdateUserHandler).Methods("PUT")

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var updatedUser user.User
	err = json.NewDecoder(rr.Body).Decode(&updatedUser)
	assert.NoError(t, err)

	assert.NotEmpty(t, updatedUser.Id)
	assert.Equal(t, updateUser.FirstName, updatedUser.FirstName)
	assert.Equal(t, updateUser.LastName, updatedUser.LastName)
	assert.Equal(t, updateUser.Nickname, updatedUser.Nickname)
	assert.Equal(t, updateUser.Email, updatedUser.Email)
	assert.Equal(t, updateUser.Country, updatedUser.Country)
	assert.NotEmpty(t, updatedUser.CreatedAt)
	assert.NotEmpty(t, updatedUser.UpdatedAt)
}
