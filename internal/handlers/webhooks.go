package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type PaymentEvent struct {
	Type       string `json:"type"`
	TransferID string `json:"transfer_id"`
	PSPID      string `json:"psp_id,omitempty"`
}

type Ack struct {
	OK         bool      `json:"ok"`
	ReceivedAt time.Time `json:"received_at"`
}

func PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var evt PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		http.Error(w, `{"code":"bad_request","message":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if evt.Type == "" || evt.TransferID == "" {
		http.Error(w, `{"code":"invalid","message":"missing required fields"}`, http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(Ack{OK: true, ReceivedAt: time.Now().UTC()})
}
