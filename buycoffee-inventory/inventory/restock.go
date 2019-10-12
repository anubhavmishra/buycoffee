package inventory

import (
	"encoding/json"
	"net/http"
)

type RestockResponse struct {
	Message string `json:"message"`
}

func RestockHandler(w http.ResponseWriter, r *http.Request) {

	response := RestockResponse{
		Message: "restocking successfull",
	}
	json.NewEncoder(w).Encode(response)
	return
}
