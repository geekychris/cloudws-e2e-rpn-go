// Package api provides the HTTP handlers for the RPN calculator service.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/geekychris/cloudws-e2e-rpn-go/internal/rpn"
)

type calculateRequest struct {
	Expression string `json:"expression"`
}

type calculateResponse struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// CalculateHandler handles POST /calculate requests.
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}

	var req calculateRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	result, err := rpn.Evaluate(req.Expression)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{Result: &result})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, calculateResponse{Error: msg})
}

func writeJSON(w http.ResponseWriter, status int, resp calculateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
