package main

import (
	"encoding/json"
	"net/http"
)

// Anything to sort out /search can go here

func (s *server) search(w http.ResponseWriter, r *http.Request) {

	var targets []string
	targets = append(targets, "Node", "GameType", "PlayerCount", "Delays")

	resp, err := json.Marshal(targets)
	if err != nil {

		http.Error(w, "cannot marshal targets response", http.StatusBadRequest)
	}
	w.Write(resp)

}
