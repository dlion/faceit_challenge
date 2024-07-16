package http

import (
	"encoding/json"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, req *http.Request) {
	//TODO: Implement healthcheck for the db
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
