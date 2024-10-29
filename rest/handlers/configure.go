package handlers

import (
	"encoding/json"
	"net/http"
	"truth/rest/server"
)

type ConfigurationHandler struct {
	configuration *server.Configuration
}

func NewConfigurationHandler(c *server.Configuration) *ConfigurationHandler {
	return &ConfigurationHandler{c}
}

func (c *ConfigurationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(c.configuration); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
