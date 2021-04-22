package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests for /acronym endpoint
func TestAcronymHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		bufferSize int
		history    History
		body       string
		want       string
		statusCode int
	}{
		{
			name:       "invalid GET method",
			method:     http.MethodGet,
			bufferSize: 1,
			history:    History{},
			body:       "",
			want:       "Method not allowed",
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "single word acronym",
			method:     http.MethodPost,
			bufferSize: 1,
			history:    History{},
			body:       `{"name":"Service"}`,
			want:       `{"name":"Service","acronym":"S"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "multiple word acronym",
			method:     http.MethodPost,
			bufferSize: 1,
			history:    History{},
			body:       `{"name":"President of the United States"}`,
			want:       `{"name":"President of the United States","acronym":"PotUS"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "invalid empty name",
			method:     http.MethodPost,
			bufferSize: 1,
			history:    History{},
			body:       `{"name":""}`,
			want:       `Name cannot be blank`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/acronym", strings.NewReader(tc.body))
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
