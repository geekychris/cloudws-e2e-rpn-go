package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandlerHappyPath(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"expression": "3 4 +"})
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	CalculateHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var resp calculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result == nil || *resp.Result != 7 {
		t.Errorf("got result %v, want 7", resp.Result)
	}
}

func TestCalculateHandlerParseError(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"expression": "1 +"})
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	CalculateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var resp calculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestCalculateHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	CalculateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
