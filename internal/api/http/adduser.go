package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
)

func (u *UserHandler) AddUserHandler(w http.ResponseWriter, req *http.Request) {
	var newUser user.NewUser
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		log.Print(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdUser, err := u.userService.NewUser(req.Context(), newUser)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Print(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
