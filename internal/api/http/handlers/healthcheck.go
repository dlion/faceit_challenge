package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type HealthCheckHandler struct {
	mongoClient *mongo.Client
}

func NewHealthCheckHandler(mongoClient *mongo.Client) *HealthCheckHandler {
	return &HealthCheckHandler{mongoClient: mongoClient}
}

func (u *HealthCheckHandler) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	err := u.mongoClient.Ping(req.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to ping the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
