package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvaluateGuessHandler_MethodNotAllowed(t *testing.T) {
	handler := NewHandler([]string{"apple"}, []string{"apple", "alley"}, true)

	req := httptest.NewRequest(http.MethodGet, "/api/guess", nil)
	rec := httptest.NewRecorder()

	handler.EvaluateGuess(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestEvaluateGuessHandler_ReturnsFeedbackAndSetsClientCookie(t *testing.T) {
	handler := NewHandler([]string{"apple"}, []string{"apple", "alley"}, true)

	body := bytes.NewBufferString(`{"guess":"alley"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/guess", body)
	rec := httptest.NewRecorder()

	handler.EvaluateGuess(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	cookies := rec.Result().Cookies()

	var clientCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "client_id" {
			clientCookie = cookie
			break
		}
	}

	if clientCookie == nil {
		t.Fatal("expected client_id cookie to be set")
	}

	if clientCookie.Value == "" {
		t.Fatal("expected client_id cookie value not to be empty")
	}

	var response struct {
		Feedback []struct {
			Letter string `json:"letter"`
			Status string `json:"status"`
		} `json:"feedback"`
		Status       string `json:"status"`
		AttemptsUsed int    `json:"attempts_used"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusInProgress {
		t.Fatalf("expected status %s, got %s", StatusInProgress, response.Status)
	}

	if response.AttemptsUsed != 1 {
		t.Fatalf("expected attempts_used 1, got %d", response.AttemptsUsed)
	}

	expected := []string{"correct", "present", "absent", "present", "absent"}

	for i, expectedStatus := range expected {
		if response.Feedback[i].Status != expectedStatus {
			t.Fatalf("index %d: expected %s, got %s", i, expectedStatus, response.Feedback[i].Status)
		}
	}
}
