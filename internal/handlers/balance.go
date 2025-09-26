package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type Balance struct {
	AccountID string    `json:"account_id"`
	Available string    `json:"available"`
	Ledger    string    `json:"ledger"`
	Currency  string    `json:"currency"`
	AsOf      time.Time `json:"as_of"`
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	resp := Balance{
		AccountID: "acct_dummy",
		Available: "0",
		Ledger:    "0",
		Currency:  "USD",
		AsOf:      time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
