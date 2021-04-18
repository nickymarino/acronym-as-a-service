package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

type AcronymRequest struct {
	Name string `json:"name":`
}

type AcronymResponse struct {
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

type History []AcronymResponse

type AcronymHandler struct {
	history *History
}

func (ah AcronymHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var ac AcronymRequest

		if err := json.NewDecoder(r.Body).Decode(&ac); err != nil {
			http.Error(w, "Can't decode body", http.StatusBadRequest)
			return
		}

		acronym := acronymFrom(ac.Name)
		response := AcronymResponse{
			Name:    ac.Name,
			Acronym: acronym,
		}

		*ah.history = append(*ah.history, response)
		json.NewEncoder(w).Encode(response)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

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
	var history History

	mux := http.NewServeMux()
	mux.Handle("/acronym", AcronymHandler{&history})
	mux.Handle("/history", HistoryHandler{&history})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
