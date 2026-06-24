package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvaluateGuessHandler_MethodNotAllowed(t *testing.T) {
	handler := NewHandler([]string{"apple"})

	req := httptest.NewRequest(http.MethodGet, "/api/guess", nil)
	rec := httptest.NewRecorder()

	handler.EvaluateGuess(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestEvaluateGuessHandler_ReturnsFeedback(t *testing.T) {
	handler := NewHandler([]string{"apple"})

	body := bytes.NewBufferString(`{"guess":"alley"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/guess", body)
	rec := httptest.NewRecorder()

	handler.EvaluateGuess(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response struct {
		Feedback []struct {
			Letter string `json:"letter"`
			Status string `json:"status"`
		} `json:"feedback"`
		IsCorrect bool `json:"is_correct"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []string{"correct", "present", "absent", "present", "absent"}

	for i, expectedStatus := range expected {
		if response.Feedback[i].Status != expectedStatus {
			t.Fatalf("index %d: expected %s, got %s", i, expectedStatus, response.Feedback[i].Status)
		}
	}
}
