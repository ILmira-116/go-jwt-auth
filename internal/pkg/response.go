package pkg

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/model"
)

type UserIDResponse struct {
	GUID string `json:"guid"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(model.ErrorResponse{Message: message})
}
