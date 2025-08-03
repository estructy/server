package health_handler

import (
	"net/http"
	"time"

	json_helper "github.com/nahtann/controlriver.com/internal/helpers/json"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond) // Simulate a small delay for processing
	json_helper.HttpResponse(w, http.StatusOK, map[string]string{
		"status":    "healthy",
		"message":   "The service is running smoothly.",
		"timestamp": "2023-10-01T12:00:00Z",
	})
}
