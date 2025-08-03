package json_helper

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorMappings struct {
	Type    int
	Message string
}

func HttpResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

func HttpError(w http.ResponseWriter, code int, message string) {
	var errType string

	switch code {
	case http.StatusNotFound:
		errType = "Not Found"
	case http.StatusBadRequest:
		errType = "Bad Request"
	case http.StatusUnauthorized:
		errType = "Unauthorized"
	case http.StatusForbidden:
		errType = "Forbidden"
	default:
		errType = "Internal Server Error"
	}

	payload := map[string]string{"error": errType, "message": message}

	HttpResponse(w, code, payload)
}

func handleError(w http.ResponseWriter, err error, errorMappings map[error]ErrorMappings,
) {
	log.Println(err)

	for knownErr, mapping := range errorMappings {
		if errors.Is(err, knownErr) {
			HttpResponse(w, mapping.Type, mapping.Message)
			return
		}
	}

	// Default case
	HttpResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
}
