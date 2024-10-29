package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"truth/rest/server"
)

type HistoryHandler struct {
	server *server.Server
}

func NewHistoryHandler(s *server.Server) *HistoryHandler {
	return &HistoryHandler{s}
}

func (s *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the request IP and body if present
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		log.Printf("Request from IP: %s, Body: %s\n", r.RemoteAddr, string(bodyBytes))
	} else {
		log.Printf("Request from IP: %s, Body: empty\n", r.RemoteAddr)
	}

	// By default send all data
	sendPrevMessages := true

	if r.Method == http.MethodPost {
		type dataBody struct {
			SendPrevMessages bool `json:"sendPrevMessages"`
		}
		var data dataBody
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		sendPrevMessages = data.SendPrevMessages
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	id, channel := s.server.AddHistoryListener(sendPrevMessages)
	defer s.server.RemoveHistoryListener(id)

	// Listen for client disconnects
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Client disconnected for HistoryListener ID: %d", id)
			return
		case message, ok := <-channel:
			if !ok {
				log.Printf("Channel closed for HistoryListener ID: %d", id)
				return
			}
			messageJson, err := json.Marshal(message)
			if err != nil {
				continue
			}
			_, err = fmt.Fprintf(w, "%s", string(messageJson))
			if err != nil {
				break
			}
			// Log the ID of the HistoryListener when a message is sent
			log.Printf("Message sent to HistoryListener ID: %d", id)
			flusher.Flush()
		}
	}
}
