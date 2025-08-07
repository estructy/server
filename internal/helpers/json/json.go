// Package jsonhelper provides utility functions for handling JSON responses and errors in HTTP requests.
package jsonhelper

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorMappings struct {
	Code    int
	Message string
}

func HTTPResponse(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

func HTTPError(w http.ResponseWriter, code int, message string) {
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

	HTTPResponse(w, code, payload)
}

func HandleError(w http.ResponseWriter, err error, errorMappings map[error]ErrorMappings,
) {
	log.Println(err)

	for knownErr, mapping := range errorMappings {
		if errors.Is(err, knownErr) {
			HTTPError(w, mapping.Code, mapping.Message)
			return
		}
	}

	// Default case
	HTTPError(w, http.StatusInternalServerError, "An unexpected error occurred")
}
