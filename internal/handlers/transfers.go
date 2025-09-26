package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type CreateTransferRequest struct {
	FromAccountID string           `json:"from_account_id"`
	ToAccountID   string           `json:"to_account_id,omitempty"`
	External      *ExternalDetails `json:"external,omitempty"`
	Amount        string           `json:"amount"`
	Currency      string           `json:"currency"`
}

type ExternalDetails struct {
	Scheme string `json:"scheme"`
	Dest   string `json:"dest"`
}

type Transfer struct {
	ID            string           `json:"id"`
	FromAccountID string           `json:"from_account_id"`
	ToAccountID   string           `json:"to_account_id,omitempty"`
	External      *ExternalDetails `json:"external,omitempty"`
	Amount        string           `json:"amount"`
	Currency      string           `json:"currency"`
	Status        string           `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
}

func CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"bad_request","message":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.FromAccountID == "" || req.Amount == "" || req.Currency == "" {
		http.Error(w, `{"code":"invalid","message":"missing required fields"}`, http.StatusUnprocessableEntity)
		return
	}
	// stubbed response
	resp := Transfer{
		ID:            "tr_" + time.Now().Format("20060102150405"),
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		External:      req.External,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "Pending",
		CreatedAt:     time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func GetTransfer(w http.ResponseWriter, r *http.Request) {
	// stubbed response always pending
	resp := Transfer{
		ID:        "tr_dummy",
		Status:    "Pending",
		CreatedAt: time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
