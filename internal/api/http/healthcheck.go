package http

import (
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	//TODO: Implement healthcheck for the db
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
