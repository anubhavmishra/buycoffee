package main

import (
	"encoding/json"
	"net/http"
)

type HealthCheckResponse struct {
	Message string `json:"message"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	response := HealthCheckResponse{
		Message: "OK",
	}
	json.NewEncoder(w).Encode(response)
	return
}
