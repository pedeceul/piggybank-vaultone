package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type CreateAccountRequest struct {
	OwnerID  string `json:"owner_id"`
	Kind     string `json:"kind"`
	Currency string `json:"currency"`
}

type Account struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Kind      string    `json:"kind"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"bad_request","message":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.OwnerID == "" || req.Kind == "" || req.Currency == "" {
		http.Error(w, `{"code":"invalid","message":"missing required fields"}`, http.StatusUnprocessableEntity)
		return
	}
	// stubbed response
	resp := Account{
		ID:        "acct_" + time.Now().Format("20060102150405"),
		OwnerID:   req.OwnerID,
		Kind:      req.Kind,
		Currency:  req.Currency,
		CreatedAt: time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
