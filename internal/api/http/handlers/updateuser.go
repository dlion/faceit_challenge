package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dlion/faceit_challenge/internal/domain/services/user"
	"github.com/gorilla/mux"
)

func (u *UserHandler) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	var updateUser user.UpdateUser
	if err := json.NewDecoder(req.Body).Decode(&updateUser); err != nil {
		log.Print(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok || id == "" {
		log.Print("Update failed, it has been provided a bad ID")
		http.Error(w, "ID parameter missing in URL", http.StatusBadRequest)
		return
	}
	updateUser.Id = id

	updatedUser, err := u.UserService.UpdateUser(req.Context(), &updateUser)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		log.Print(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
