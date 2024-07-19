package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (u *UserHandler) RemoveUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok || id == "" {
		log.Print("Delete failed, it has been provided a bad ID")
		http.Error(w, "ID parameter missing in URL", http.StatusBadRequest)
		return
	}

	err := u.UserService.RemoveUser(req.Context(), id)
	if err != nil {
		log.Print("Delete failed, ", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
}
