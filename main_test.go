package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAcronymHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		bufferSize int
		history    History
		want       string
		statusCode int
	}{
		{
			name:       "Invalid GET method",
			method:     http.MethodGet,
			bufferSize: 1,
			history:    History{},
			want:       "Method not allowed",
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/acronym", nil)
			responseRecorder := httptest.NewRecorder()

			AcronymHandler{&tc.bufferSize, &tc.history}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
