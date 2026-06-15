// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
)

// healthResponse is the JSON body returned by the health endpoint.
type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// HealthHandler returns an http.HandlerFunc that responds with a 200 OK
// JSON payload indicating the service is healthy.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(healthResponse{
			Status:  "ok",
			Service: "daedalus",
		})
	}
}
