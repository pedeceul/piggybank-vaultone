package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaymentWebhook_Validation(t *testing.T) {
	body := bytes.NewBufferString(`{"type":"PaymentSettled","transfer_id":"tr1"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/webhooks/payment_event", body)
	rr := httptest.NewRecorder()
	PaymentWebhook(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rr.Code)
	}
}

func TestPaymentWebhook_Missing(t *testing.T) {
	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/webhooks/payment_event", body)
	rr := httptest.NewRecorder()
	PaymentWebhook(rr, req)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
}
