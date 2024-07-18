package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRemoveUserHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/user/66981a71a4fd0f7ff33251b1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockedUserService := new(MockUserService)
	userHandler := UserHandler{UserService: mockedUserService}
	mockedUserService.On("RemoveUser").Return(nil)
	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", userHandler.RemoveUserHandler).Methods("DELETE")

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
