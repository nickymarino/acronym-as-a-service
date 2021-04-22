package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

// Payload for POSTing to /acronym
type AcronymRequest struct {
	Name string `json:"name"`
}

// Response from POST tp /acronym
type AcronymResponse struct {
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

// Hold N previous acronyms
type History []AcronymResponse

type AcronymHandler struct {
	bufferSize *int
	history    *History
}

// Handler for /acronym
func (ah AcronymHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var ac AcronymRequest

		if err := json.NewDecoder(r.Body).Decode(&ac); err != nil {
			// Fail if JSON is bad
			http.Error(w, "Can't decode body", http.StatusBadRequest)
			return
		}

		if len(ac.Name) == 0 {
			// Fail if name is empty
			http.Error(w, "Name cannot be blank", http.StatusBadRequest)
			return
		}

		// Create the acronym from the name
		acronym := acronymFrom(ac.Name)
		response := AcronymResponse{
			Name:    ac.Name,
			Acronym: acronym,
		}

		// Record the response for history and return it
		ah.Record(response)
		json.NewEncoder(w).Encode(response)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Add the new acronym to the history, and rotate out older histories
// if the history is full
func (ah AcronymHandler) Record(ar AcronymResponse) {
	if len(*ah.history) >= *ah.bufferSize {
		// Keep 1 slot open for the new item to add
		startIndex := len(*ah.history) - *ah.bufferSize + 1
		*ah.history = (*ah.history)[startIndex:]
	}
	*ah.history = append(*ah.history, ar)
}

// Generate an acronym from a name
func acronymFrom(name string) string {
	var acronymRunes []rune

	for _, word := range strings.Split(name, " ") {
		// Convert word into bytes to grab the first rune of the string
		runes := []byte(word)
		r, _ := utf8.DecodeRune(runes)
		acronymRunes = append(acronymRunes, r)
	}

	acronym := string(acronymRunes)
	return acronym
}

type HistoryHandler struct {
	history *History
}

// Handler for /history endpoint
func (hh HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if len(*hh.history) == 0 {
			http.Error(w, "Error: No acronym history found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(hh.history)
	default:
		json.NewEncoder(w).Encode(hh.history)
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	bufferSize := 20
	var history History

	mux := http.NewServeMux()
	mux.Handle("/acronym", AcronymHandler{&bufferSize, &history})
	mux.Handle("/history", HistoryHandler{&history})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
