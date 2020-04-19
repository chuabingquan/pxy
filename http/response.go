package http

import (
	"encoding/json"
	"net/http"
)

type standardResponse struct {
	Message string `json:"message"`
}

func writeJSONMessage(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
