package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/users?limit=10&offset=0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockedUserService := new(MockUserService)
	userHandler := UserHandler{UserService: mockedUserService}
	var mockedUsers []*user.User
	mockedUsers = append(mockedUsers, &user.User{
		Id:        "66981a71a4fd0f7ff33251b1",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "USA",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}, &user.User{
		Id:        "66981a71a4fd0f7ff33251b2",
		FirstName: "Pippo",
		LastName:  "Pluto",
		Nickname:  "Pipp2",
		Email:     "pipp.plut@example.com",
		Country:   "UK",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	})

	mockedUserService.On("GetUsers").Return(mockedUsers)
	router := mux.NewRouter()
	router.HandleFunc("/api/users", userHandler.GetUsersHandler).Methods("GET")

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
